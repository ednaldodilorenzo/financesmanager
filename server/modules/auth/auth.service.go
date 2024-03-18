package auth

import (
	"errors"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	ExecuteAuthentication(username string, password string) (*model.User, error)
	RegisterUser(user *model.User) error
}

type AuthServiceStruct struct {
	repository AuthRepository
}

func NewAuthService() AuthService {
	return &AuthServiceStruct{
		repository: NewAuthRepository(),
	}
}

func (a *AuthServiceStruct) ExecuteAuthentication(username string, password string) (*model.User, error) {

	user, err := a.repository.FindUserByEmail(username)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("password does not match")
	}

	return user, nil
}

func (a *AuthServiceStruct) RegisterUser(user *model.User) error {
	return config.TxWrapper(func() error {
		currentUser, err := a.repository.FindUserByEmail(user.Email)

		if err != nil {
			return err
		}

		if currentUser != nil {
			return errors.New("duplicate key value violates unique")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

		if err != nil {
			return err
		}

		user.Password = string(hashedPassword)

		err = a.repository.Create(user)

		if err != nil {
			return err
		}

		return nil
	})
}