package repository

import (
	"errors"
	
	"go.uber.org/dig"
	"gorm.io/gorm"
	"twitter/internal/infrastructure/mysql"
	"twitter/internal/models"
)

type TweetsRepositoryParams struct {
	dig.In

	MySQLClientHandler *mysql.MySQLClientHandler
}

type TweetsRepository struct {
	db *gorm.DB
}

func NewTweetsRepository(params TweetsRepositoryParams) *TweetsRepository {
	return &TweetsRepository{db: params.MySQLClientHandler.DB}
}

func (r *TweetsRepository) PostTweets(userID int, data string) (*models.Tweet, error) {
	// Create a tweet
	tweet := &models.Tweet{
		UserID:      userID, // Foreign key reference
		Description: data,
	}

	if err := r.db.Create(&tweet).Error; err != nil {
		return nil, err
	} else {
		return tweet, nil
	}
}

func (r *TweetsRepository) GetUserTweets(userID int, sortingType string) ([]*models.Tweet, error) {
	var feeds []*models.Tweet
	query := r.db.Table("tweets").Where("user_id = ?", userID)
	if sortingType == "desc" {
		query = query.Order("created_at desc")
	} else {
		query = query.Order("created_at asc")
	}
	// get all tweets of those users
	if err := query.Find(&feeds).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return feeds, nil
}
