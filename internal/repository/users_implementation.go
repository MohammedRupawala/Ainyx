package repository

import (
	"context"

	"github.com/Ainyx-backend/db/sqlc"
	"github.com/Ainyx-backend/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
)

type userRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(q *sqlc.Queries) UserRepository {
	return &userRepository{
		queries: q,
	}
}

func (r *userRepository) CreateUser(
	ctx context.Context,
	data models.CreateUserInput,
) (models.User, error) {

	params := sqlc.CreateUserParams{
		Name: data.Name,
		Dob: pgtype.Date{
			Time:  data.DOB,
			Valid: true,
		},
	}
	response, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return models.User{}, err
	}

	return mapUser(response), nil
}

func (r *userRepository) GetByID(
	ctx context.Context,
	id int32,
) (models.User, error) {

	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}

	return mapUser(user), nil
}

func (r *userRepository) GetAll(
	ctx context.Context,
	limit,
	offset int32,
) ([]models.User, error) {

	users, err := r.queries.ListUsers(
		ctx,
		sqlc.ListUsersParams{
			Limit:  limit,
			Offset: offset,
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]models.User, 0, len(users))
	for _, user := range users {
		result = append(result, mapUser(user))
	}

	return result, nil
}

func (r *userRepository) Update(
	ctx context.Context,
	params models.UpdateUserInput,
) (models.User, error) {
	updated, err := r.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:   params.ID,
		Name: params.Name,
		Dob: pgtype.Date{
			Time:  params.DOB,
			Valid: true,
		},
	})
	if err != nil {
		return models.User{}, err
	}

	return mapUser(updated), nil
}

func (r *userRepository) Delete(
	ctx context.Context,
	id int32,
) error {

	return r.queries.DeleteUser(ctx, id)
}

func mapUser(user sqlc.User) models.User {
	result := models.User{
		ID:   user.ID,
		Name: user.Name,
	}

	if user.Dob.Valid {
		result.DOB = user.Dob.Time
	}

	return result
}
