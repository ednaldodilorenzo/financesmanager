package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/dto"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/util"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	ExecuteAuthentication(username string, password string) (*SignInResponse, error)
	RegisterUserOAuthUser(userDto dto.UserOAuthRegistrationRegistration) (*SignInResponse, error)
	StartRegistrationProcess(string) error
	RegisterUserWithToken(*SignUpInput) error
	ChangePassword(int, *ChangePasswordRequest) error
	StartRecoverPasswordProcess(string) error
	RedefinePassword(*RedefinePasswordRequest) error
}

type AuthServiceStruct struct {
	repository  AuthRepository
	emailSender util.EmailSender
	settings    *config.Settings
}

func NewAuthService(authRepository AuthRepository, emailSender util.EmailSender, settings *config.Settings) AuthService {
	return &AuthServiceStruct{
		repository:  authRepository,
		emailSender: emailSender,
		settings:    settings,
	}
}

func (a *AuthServiceStruct) ExecuteAuthentication(username string, password string) (*SignInResponse, error) {

	user, err := a.repository.FindUserByEmail(username)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, util.NewAPIError(util.ErrBusiness, []string{"invalid username or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
		return nil, util.NewAPIError(util.ErrBusiness, []string{"password does not match"})
	}

	expTime := 30 * time.Minute
	tokenString, err := util.GenerateToken(user.ID, &a.settings.AppSettings.JwtKey, &expTime)

	if err != nil {
		return nil, err
	}

	response := &SignInResponse{
		ID:    user.ID,
		Name:  user.Name,
		Token: *tokenString,
	}

	return response, nil
}

func (a *AuthServiceStruct) RegisterUserOAuthUser(userDto dto.UserOAuthRegistrationRegistration) (*SignInResponse, error) {

	currentUser, err := a.repository.FindUserByEmail(userDto.Email)

	if err != nil {
		return nil, err
	}

	if currentUser == nil {
		now := time.Now().UTC()

		currentUser = &model.User{
			Name:      userDto.Name,
			Email:     userDto.Email,
			CreatedAt: &now,
		}

		currentUser, err = a.repository.Create(currentUser)

		if err != nil {
			return nil, err
		}
	}

	expirationTime := 30 * time.Minute
	tokenString, err := util.GenerateToken(currentUser.ID, &a.settings.AppSettings.JwtKey, &expirationTime)

	if err != nil {
		return nil, err
	}

	result := &SignInResponse{
		ID:    currentUser.ID,
		Token: *tokenString,
		Name:  currentUser.Name,
	}

	return result, nil
}

func (a *AuthServiceStruct) StartRegistrationProcess(email string) error {
	user, err := a.repository.FindUserByEmail(email)

	if err != nil {
		return err
	}

	if user != nil {
		return util.NewAPIError(util.ErrBusiness, []string{"User already registered"})
	}

	expirationTime := 30 * time.Minute
	tokenString, err := util.GenerateToken(&email, &a.settings.AppSettings.JwtKey, &expirationTime)

	if err != nil {
		return err
	}

	err = a.emailSender.SendEmail(email, "Confirmação de email", fmt.Sprintf("Clique no link abaixo para confirmar seu email<br><br><a href=\"%s\">Clique aqui.</a>", fmt.Sprintf("%s/register/%s", a.settings.AppSettings.Url, *tokenString)))

	if err != nil {
		return err
	}

	return nil
}

func (a *AuthServiceStruct) RegisterUserWithToken(signin *SignUpInput) error {
	email, err := util.ExtractSubContent(&signin.Token, &a.settings.AppSettings.JwtKey)

	if err != nil {
		return err
	}

	user, err := a.repository.FindUserByEmail(*email)

	if err != nil {
		return err
	}

	if user != nil {
		return util.NewAPIError(util.ErrBusiness, []string{"User already registered"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signin.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	now := time.Now()
	password := string(hashedPassword)
	newUser := &model.User{
		Name:      signin.Name,
		Email:     *email,
		Password:  &password,
		CreatedAt: &now,
	}

	if _, err = a.repository.Create(newUser); err != nil {
		return errors.New("failed to register new user")
	}

	return nil
}

func (a *AuthServiceStruct) ChangePassword(userId int, changePassword *ChangePasswordRequest) error {
	user, err := a.repository.FindById(userId)

	if err != nil {
		return err
	}

	if user == nil {
		return util.NewAPIError(util.ErrNotFound, []string{"User not Found"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(changePassword.Password)); err != nil {
		return util.NewAPIError(util.ErrBusiness, []string{"Senha atual não corresponde a cadastada"})
	}

	if changePassword.NewPassword != changePassword.CofirmNewPassword {
		return util.NewAPIError(util.ErrBusiness, []string{"Nova senha não confere com a confirmação"})
	}

	if changePassword.NewPassword == changePassword.Password {
		return util.NewAPIError(util.ErrBusiness, []string{"Nova senha não deve ser igual a anterior"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(changePassword.NewPassword), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	password := string(hashedPassword)
	user.Password = &password

	err = a.repository.Update(userId, user)

	if err != nil {
		return err
	}

	return nil
}

func (a *AuthServiceStruct) StartRecoverPasswordProcess(email string) error {
	user, err := a.repository.FindUserByEmail(email)

	if err != nil {
		return err
	}

	if user == nil {
		return util.NewAPIError(util.ErrNotFound, []string{"Usuário não registrado"})
	}

	expirationTime := 30 * time.Minute
	tokenString, err := util.GenerateToken(&email, &a.settings.AppSettings.JwtKey, &expirationTime)

	if err != nil {
		return err
	}

	err = a.emailSender.SendEmail(email, "Recuperação de Conta", fmt.Sprintf("Clique no link abaixo para registrar uma nova senha.<br><br><a href=\"%s\">Clique aqui.</a>", fmt.Sprintf("%s/redefine/%s", a.settings.AppSettings.Url, *tokenString)))

	if err != nil {
		return err
	}

	return nil
}

func (a *AuthServiceStruct) RedefinePassword(redefinePassword *RedefinePasswordRequest) error {

	email, err := util.ExtractSubContent(&redefinePassword.Token, &a.settings.AppSettings.JwtKey)

	if err != nil {
		return err
	}

	user, err := a.repository.FindUserByEmail(*email)

	if err != nil {
		return err
	}

	if user == nil {
		return util.NewAPIError(util.ErrNotFound, []string{"Usuário não encontrado"})
	}

	if redefinePassword.Password != redefinePassword.CofirmNewPassword {
		return util.NewAPIError(util.ErrBusiness, []string{"Senha diferente da confirmação"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(redefinePassword.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	password := string(hashedPassword)
	user.Password = &password

	if err := a.repository.Update(int(user.ID), user); err != nil {
		return err
	}

	return nil
}
