package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/util"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	ExecuteAuthentication(username string, password string) (*model.User, error)
	RegisterUser(user *model.User) error
	StartRegistrationProcess(string) error
	RegisterUserWithToken(*SignUpInput) error
	ChangePassword(int, string) error
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

func (a *AuthServiceStruct) ExecuteAuthentication(username string, password string) (*model.User, error) {

	user, err := a.repository.FindUserByEmail(username)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, util.NewBusinessError("password does not match", err, util.BE_PASSWORD_DO_NOT_MATCH)
	}

	return user, nil
}

func (a *AuthServiceStruct) RegisterUser(user *model.User) error {

	currentUser, err := a.repository.FindUserByEmail(user.Email)

	if err != nil {
		return err
	}

	if currentUser != nil {
		return util.NewBusinessError("User already registered!", nil, util.BE_USER_ALREADY_REGISTERED)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	now := time.Now().UTC()
	user.CreatedAt = &now

	err = a.repository.Create(user)

	if err != nil {
		return err
	}

	tokenByte := jwt.New(jwt.SigningMethodHS256)

	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["sub"] = user.Email
	claims["exp"] = now.Add(30 * time.Minute).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(a.settings.AppSettings.JwtKey))

	if err != nil {
		return err
	}

	err = a.emailSender.SendEmail(user.Email, "Confirmação de email", fmt.Sprintf("Clique no link abaixo para confirmar seu email<br><br><a href=\"%s\">Clique aqui.</a>", fmt.Sprintf("%s/verify/%s", a.settings.AppSettings.Url, tokenString)))

	if err != nil {
		return err
	}

	return nil
}

func (a *AuthServiceStruct) StartRegistrationProcess(email string) error {
	user, err := a.repository.FindUserByEmail(email)

	if err != nil {
		return err
	}

	if user != nil {
		return util.NewBusinessError("User already registered", nil, util.BE_USER_ALREADY_REGISTERED)
	}

	tokenByte := jwt.New(jwt.SigningMethodHS256)

	claims := tokenByte.Claims.(jwt.MapClaims)

	now := time.Now().UTC()

	claims["sub"] = email
	claims["exp"] = now.Add(30 * time.Minute).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(a.settings.AppSettings.JwtKey))

	if err != nil {
		return err
	}

	err = a.emailSender.SendEmail(email, "Confirmação de email", fmt.Sprintf("Clique no link abaixo para confirmar seu email<br><br><a href=\"%s\">Clique aqui.</a>", fmt.Sprintf("%s/register/%s", a.settings.AppSettings.Url, tokenString)))

	if err != nil {
		return err
	}

	return nil
}

func (a *AuthServiceStruct) RegisterUserWithToken(signin *SignUpInput) error {
	tokenByte, err := jwt.Parse(signin.Token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}

		return []byte(a.settings.AppSettings.JwtKey), nil
	})

	if err != nil {
		return err
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return errors.New("invalid token")

	}

	email := fmt.Sprint(claims["sub"])

	user, err := a.repository.FindUserByEmail(email)

	if err != nil {
		return err
	}

	if user != nil {
		return util.NewBusinessError("User already registered", nil, util.BE_USER_ALREADY_REGISTERED)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signin.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	now := time.Now()

	newUser := &model.User{
		Name:      signin.Name,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: &now,
	}

	if err = a.repository.Create(newUser); err != nil {
		return errors.New("failed to register new user")
	}

	return nil
}

func (a *AuthServiceStruct) ChangePassword(userId int, newPassword string) error {
	user, err := a.repository.FindById(userId)

	if err != nil {
		return err
	}

	if user == nil {
		return util.NewBusinessError("User not Found", nil, util.BE_INPUT_VALIDATION_ERROR)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	err = a.repository.Update(userId, user)

	if err != nil {
		return err
	}

	return nil
}
