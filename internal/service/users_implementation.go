package service

import (
	"context"
	"time"

	"github.com/Ainyx-backend/internal/models"
	"github.com/Ainyx-backend/internal/repository"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(
	ctx context.Context,
	userParams models.CreateUserInput,
) (models.User, error) {
	user, err := s.repo.CreateUser(ctx, userParams)
	if err != nil {
		return models.User{}, err
	}

	return withAge(user), nil
}

func (s *userService) GetByID(
	ctx context.Context,
	id int32,
) (models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}

	return withAge(user), nil
}

func (s *userService) GetAll(
	ctx context.Context,
	limit,
	offset int32,
) ([]models.User, error) {
	users, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]models.User, 0, len(users))
	for _, user := range users {
		result = append(result, withAge(user))
	}

	return result, nil
}

func (s *userService) Update(
	ctx context.Context,
	params models.UpdateUserInput,
) (models.User, error) {
	user, err := s.repo.Update(ctx, params)
	if err != nil {
		return models.User{}, err
	}

	return withAge(user), nil
}

func (s *userService) Delete(ctx context.Context, id int32) error {

	return s.repo.Delete(ctx, id)
}

func withAge(user models.User) models.User {
	user.Age = calculateAge(user.DOB)
	return user
}

func calculateAge(dob time.Time) int {
	if dob.IsZero() {
		return 0
	}

	today := time.Now().UTC()
	birth := dob.UTC()
	age := today.Year() - birth.Year()
	if today.Month() < birth.Month() || (today.Month() == birth.Month() && today.Day() < birth.Day()) {
		age--
	}

	if age < 0 {
		return 0
	}
	return age
}
