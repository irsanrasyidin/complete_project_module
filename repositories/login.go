package repositories

import (
	"github.com/irsanrasyidin/complete_project_module/models"
	"github.com/irsanrasyidin/complete_project_module/utils"

	"gorm.io/gorm"
)

type LoginRepo interface {
	LoginUser(user models.UserCredential) (string, error)
}

type LoginRepoImpl struct {
	db *gorm.DB
}

func (lgnRepo *LoginRepoImpl) LoginUser(user models.UserCredential) (string, error) {
	if err := lgnRepo.db.Where("username = ? AND password = ?", user.Username, user.Password).First(user.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", &utils.AppError{
				ErrorCode:    101,
				ErrorMessage: "Invalid username or password",
			}
		}
		return "", &utils.AppError{
			ErrorCode:    102,
			ErrorMessage: err.Error(),
		}
	}
	return user.ID, nil
}

func NewLoginRepo() LoginRepo {
	return &LoginRepoImpl{}
}
