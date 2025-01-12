package repository

import (
	"errors"

	"go.uber.org/dig"
	"gorm.io/gorm"
	"twitter/internal/infrastructure/mysql"
	"twitter/internal/models"
)

type UserRepositoryParams struct {
	dig.In

	MySQLClientHandler *mysql.MySQLClientHandler
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(params UserRepositoryParams) *UserRepository {
	return &UserRepository{db: params.MySQLClientHandler.DB}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FollowUserByID(followerID int, userID int) (*models.Followers, error) {
	follow := &models.Followers{
		FollowingUserID: followerID,
		UserID:          userID,
	}
	if err := r.db.Create(follow).Error; err != nil {
		return nil, err
	}
	return follow, nil
}
