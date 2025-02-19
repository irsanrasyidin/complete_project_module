//go:generate mockgen -source=user.go -destination=mocks/user_mock.go UserRepo

package repositories

import (
	"github.com/irsanrasyidin/complete_project_module/models"
	"github.com/irsanrasyidin/complete_project_module/utils"

	"gorm.io/gorm"
)

type UserRepo interface {
	InsertUser(lgc *models.UserModel, password string) error
	GetUserByID(ID string) (*models.UserModel, error)
	UpdateUserByUsername(lgc *models.UserModel) error
	DeleteUserById(id string) error
}

type UserRepoImpl struct {
	db *gorm.DB
}

func (usrRepo *UserRepoImpl) InsertUser(lgc *models.UserModel, password string) error {
	tx := usrRepo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(lgc).Error; err != nil {
		tx.Rollback()
		return &utils.AppError{
			ErrorCode:    111,
			ErrorMessage: err.Error(),
		}
	}

	credential := &models.UserCredential{
		ID:       lgc.ID,
		Username: lgc.Username,
		Password: password,
	}

	if err := tx.Create(credential).Error; err != nil {
		tx.Rollback()
		return &utils.AppError{
			ErrorCode:    112,
			ErrorMessage: err.Error(),
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return &utils.AppError{
			ErrorCode:    113,
			ErrorMessage: err.Error(),
		}
	}
	return nil
}

func (usrRepo *UserRepoImpl) GetUserByID(ID string) (*models.UserModel, error) {
	var users *models.UserModel
	if err := usrRepo.db.Where("id = ?", ID).Find(&users).Error; err != nil {
		return nil, &utils.AppError{
			ErrorCode:    112,
			ErrorMessage: err.Error(),
		}
	}
	return users, nil
}

func (usrRepo *UserRepoImpl) UpdateUserByUsername(lgc *models.UserModel) error {
	if err := usrRepo.db.Model(&models.UserModel{}).Where("username = ?", lgc.Username).Updates(lgc).Error; err != nil {
		return &utils.AppError{
			ErrorCode:    113,
			ErrorMessage: err.Error(),
		}
	}
	return nil
}

func (usrRepo *UserRepoImpl) DeleteUserById(id string) error {
	if err := usrRepo.db.Where("id = ?", id).Delete(&models.UserModel{}).Error; err != nil {
		return &utils.AppError{
			ErrorCode:    114,
			ErrorMessage: err.Error(),
		}
	}
	return nil
}

func NewUserRepo() UserRepo {
	return &UserRepoImpl{}
}
