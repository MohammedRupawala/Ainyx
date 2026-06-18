package repository

import (
	"context"

	"github.com/Ainyx-backend/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, userParams models.CreateUserInput) (models.User, error)
	GetByID(ctx context.Context, id int32) (models.User, error)
	GetAll(ctx context.Context, limit, offset int32) ([]models.User, error)
	Update(ctx context.Context, params models.UpdateUserInput) (models.User, error)
	Delete(ctx context.Context, id int32) error
}