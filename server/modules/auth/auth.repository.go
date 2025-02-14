package auth

import (
	"errors"
	"strings"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Create(user *model.User) error
	FindUserByEmail(email string) (*model.User, error)
	FindById(int) (*model.User, error)
	Update(id int, item *model.User) error
}

type AuthRepositoryStruct struct {
	dbConfig *config.Database
}

func NewAuthRepository(database *config.Database) AuthRepository {
	return &AuthRepositoryStruct{
		dbConfig: database,
	}
}

func (a *AuthRepositoryStruct) FindUserByEmail(email string) (*model.User, error) {
	var user model.User

	if err := a.dbConfig.DB.First(&user, "email = ?", strings.ToLower(email)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (a *AuthRepositoryStruct) FindById(id int) (*model.User, error) {
	var user model.User

	if err := a.dbConfig.DB.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (a *AuthRepositoryStruct) Create(user *model.User) error {
	result := a.dbConfig.DB.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (g *AuthRepositoryStruct) Update(id int, item *model.User) error {
	if err := g.dbConfig.DB.Model(&item).Where("id = ?", id).Updates(item).Error; err != nil {
		return err
	}

	return nil
}
