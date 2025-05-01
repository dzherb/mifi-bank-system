package auth

import (
	"fmt"

	"github.com/dzherb/mifi-bank-system/internal/models"
	repo "github.com/dzherb/mifi-bank-system/internal/repository"
	"github.com/dzherb/mifi-bank-system/internal/security"
)

type Service interface {
	Login(email, password string) (AccessPayload, error)
	Register(email, username, password string) (AccessPayload, error)
}

type AccessPayload struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func NewService() Service {
	return &ServiceImpl{
		emailVal: DefaultEmailValidator,
		userVal:  DefaultUsernameValidator,
		passVal:  DefaultPasswordValidator,
	}
}

type ServiceImpl struct {
	emailVal Validator
	userVal  Validator
	passVal  Validator
}

func (s *ServiceImpl) Login(email, password string) (AccessPayload, error) {
	const errFmt = "failed to login: %w"

	// Validate email and password early to avoid unnecessary DB access
	err := s.validateLoginInput(email, password)
	if err != nil {
		return AccessPayload{}, fmt.Errorf(errFmt, err)
	}

	ur := repo.NewUserRepository()

	user, err := ur.GetByCredentials(email, password)
	if err != nil {
		return AccessPayload{}, fmt.Errorf(errFmt, err)
	}

	return s.issueToken(user, errFmt)
}

func (s *ServiceImpl) Register(
	email,
	username,
	password string,
) (AccessPayload, error) {
	const errFmt = "failed to register: %w"

	err := s.validateRegisterInput(email, username, password)
	if err != nil {
		return AccessPayload{}, fmt.Errorf(errFmt, err)
	}

	ur := repo.NewUserRepository()

	user, err := ur.Create(models.User{
		Email:    email,
		Username: username,
		Password: password,
	})

	if err != nil {
		return AccessPayload{}, fmt.Errorf(errFmt, err)
	}

	return s.issueToken(user, errFmt)
}

func (s *ServiceImpl) validateLoginInput(email, password string) error {
	if err := s.emailVal.Validate(email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	if err := s.passVal.Validate(password); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	return nil
}

func (s *ServiceImpl) validateRegisterInput(
	email,
	username,
	password string,
) error {
	if err := s.emailVal.Validate(email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	if err := s.userVal.Validate(username); err != nil {
		return fmt.Errorf("invalid username: %w", err)
	}

	if err := s.passVal.Validate(password); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	return nil
}

func (s *ServiceImpl) issueToken(
	user models.User,
	context string,
) (AccessPayload, error) {
	token, err := security.IssueAccessToken(user.ID)
	if err != nil {
		return AccessPayload{}, fmt.Errorf(context, err)
	}

	return AccessPayload{Token: token, User: user}, nil
}
