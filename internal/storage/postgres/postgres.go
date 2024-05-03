package postgres

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main/internal/domain/models"
	"main/internal/storage"
	"time"
)

type Storage struct {
	DB *gorm.DB
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"ssl_mode"`
}

// New create a connection to database and returns a structure pointer to the created database
func New(p Postgres) (*Storage, error) {
	const op = "internal.storage.postgres.New"

	// Data Source Name
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		p.Host, p.Port, p.User, p.DBName, p.Password, p.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return &Storage{}, err
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return &Storage{}, err
	}

	return &Storage{
		DB: db,
	}, nil
}

func (db *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	done := make(chan struct{}, 1)
	errDone := make(chan struct{}, 1)
	var id int64

	go func() {
		var user models.User
		checkUserExists := db.DB.WithContext(ctx).Find(&user, "email = ?", email)
		if checkUserExists.RowsAffected != 0 {
			errDone <- struct{}{}
		} else {
			var users []models.User
			ids := db.DB.WithContext(ctx).Find(&users)
			userID := ids.RowsAffected
			db.DB.Create(&models.User{
				ID:       userID + 1,
				Email:    email,
				PassHash: passHash,
			})

			id = userID
			done <- struct{}{}
		}
		close(errDone)
		close(done)
	}()

	select {
	case <-errDone:
		return 0, storage.ErrUserExists
	case <-ctx.Done():
		return 0, storage.ErrConnection
	case <-done:
		return id + 1, nil
	}
}

func (db *Storage) CreateUser(ctx context.Context, email string) (models.User, error) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	done := make(chan struct{}, 1)
	errDone := make(chan struct{}, 1)
	var user models.User

	go func() {
		db.DB.WithContext(ctx).Where(&models.User{Email: email}).First(&user)
		if user.ID == 0 {
			errDone <- struct{}{}
		} else {
			user = models.User{
				ID:       user.ID,
				Email:    user.Email,
				PassHash: user.PassHash,
			}
			done <- struct{}{}
		}
		close(errDone)
		close(done)
	}()

	select {
	case <-errDone:
		return models.User{}, storage.ErrUserNotFound
	case <-ctx.Done():
		return models.User{}, storage.ErrConnection
	case <-done:
		return user, nil
	}
}
