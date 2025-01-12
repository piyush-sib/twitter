package repository

import (
	"errors"
	"go.uber.org/dig"
	"gorm.io/gorm"
	"log"

	"twitter/internal/infrastructure/mysql"
	"twitter/internal/models"
)

type FeedsRepositoryParams struct {
	dig.In

	MySQLClientHandler *mysql.MySQLClientHandler
}

type FeedsRepository struct {
	db *gorm.DB
}

func NewFeedsRepository(params FeedsRepositoryParams) *FeedsRepository {
	return &FeedsRepository{db: params.MySQLClientHandler.DB}
}

func (r *FeedsRepository) GetUserFeeds(userID int, sortingType string) ([]*models.Tweet, error) {
	var feeds []*models.Tweet
	var followers []*models.Followers

	var followersIDs []int
	// get all followings users first
	r.db.Table("followers").Select("following_user_id").Where("user_id = ?", userID).Find(&followers)
	log.Printf("this is the followers %v", followers)

	for _, v := range followers {
		followersIDs = append(followersIDs, v.FollowingUserID)
	}
	log.Printf("this is the followersIDs %v", followersIDs)

	// get all tweets of those users
	if err := r.db.Table("tweets").Where("user_id IN ?", followersIDs).Find(&feeds).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return feeds, nil
}
