package users

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"testapi/internal/entity"
)

type Repository interface {
	Get(ctx context.Context, id int64) (entity.User, error)

	Count(ctx context.Context) (int, error)

	Query(ctx context.Context, offset, limit int) ([]entity.User, error)

	Create(ctx context.Context, req entity.User) (entity.User, error)

	Update(ctx context.Context, id int64, req entity.User) (entity.User, error)

	Delete(ctx context.Context, id int64) error
}

type repository struct {
	dbClient *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) Repository {
	return repository{dbClient: db}
}

func (r repository) Get(ctx context.Context, id int64) (entity.User, error) {
	var user = entity.User{}

	sql, args, err := sq.Select("*").
		From("users").
		Where(sq.Eq{"user_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	err = pgxscan.Get(ctx, r.dbClient, &user, sql, args...)
	return user, err
}

func (r repository) Count(ctx context.Context) (int, error) {
	panic("implement me")
}

func (r repository) Query(ctx context.Context, offset, limit int) ([]entity.User, error) {
	var users []entity.User

	sql, _, err := sq.Select("*").
		From("users").
		ToSql()

	err = pgxscan.Select(ctx, r.dbClient, &users, sql)

	return users, err

}

func (r repository) Create(ctx context.Context, req entity.User) (entity.User, error) {
	var newUser = entity.User{}
	var existingUser = entity.User{}

	checkUserQuery, args, err :=
		sq.Select("*").
			From("users").
			Where(sq.Eq{"email": req.Email}).
			PlaceholderFormat(sq.Dollar).
			ToSql()

	err = pgxscan.Get(ctx, r.dbClient, &existingUser, checkUserQuery)
	if err == nil {
		return entity.User{}, fmt.Errorf("%v is already in use", existingUser.Email)
	}

	insertQuery, args, err := sq.Insert("users").
		Columns("first_name",
			"last_name",
			"address",
			"email",
			"password").
		Suffix("RETURNING user_id").
		Values(req.FirstName, req.LastName, req.Address, req.Email, req.Password).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	fmt.Println(insertQuery)

	err = r.dbClient.QueryRow(ctx, insertQuery, args...).Scan(&newUser.ID)

	if err != nil {
		return entity.User{}, err
	}

	newUser.FirstName = req.FirstName
	newUser.LastName = req.LastName
	newUser.Address = req.Address
	newUser.Email = req.Email

	return newUser, nil
}

func (r repository) Update(ctx context.Context, id int64, req entity.User) (entity.User, error) {
	var updatedUser = entity.User{}

	q, args, err := sq.Update("users").
		Set("first_name", req.FirstName).
		Set("last_name", req.LastName).
		Set("address", req.Address).
		Set("email", req.Email).
		Set("password", req.Password).
		Where(sq.Eq{"user_id": id}).Suffix("RETURNING user_id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	err = r.dbClient.QueryRow(ctx, q, args...).Scan(&updatedUser.ID)

	if err != nil {
		return entity.User{}, err
	}

	updatedUser.ID = id
	updatedUser.FirstName = req.FirstName
	updatedUser.LastName = req.LastName
	updatedUser.Address = req.Address
	updatedUser.Email = req.Email

	return updatedUser, nil
}

func (r repository) Delete(ctx context.Context, id int64) error {

	q, args, err := sq.Delete("users").
		Where(sq.Eq{"user_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	_, err = r.dbClient.Exec(ctx, q, args...)

	if err != nil {
		return err
	}

	return nil
}
