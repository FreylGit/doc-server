package user

import (
	"context"
	"doc-server/internal/models"
	"doc-server/internal/storage"
	"doc-server/internal/storage/db"
	modelsRepo "doc-server/internal/storage/db/pg/models"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	tableName          = "users"
	idColumn           = "id"
	loginColumn        = "login"
	passwordHashColumn = "password_hash"
	createdAtColumn    = "created_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) storage.UserRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, user models.User) error {
	builder := squirrel.Insert(tableName).
		Columns(loginColumn, passwordHashColumn).
		Values(user.Login, user.Password).
		PlaceholderFormat(squirrel.Dollar)
	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	tag, err := r.db.Exec(ctx, query, args...)
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("error: not create")
	}
	return nil
}

func (r *repo) Get(ctx context.Context, login string) (models.User, error) {
	builder := squirrel.Select(idColumn, loginColumn, passwordHashColumn, createdAtColumn).
		From(tableName).
		Where(squirrel.Eq{loginColumn: login}).
		Limit(1).
		PlaceholderFormat(squirrel.Dollar)
	query, args, err := builder.ToSql()
	if err != nil {
		return models.User{}, err
	}
	row := r.db.QueryRow(ctx, query, args...)
	var user modelsRepo.User
	err = row.Scan(&user.Id, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil {
		return models.User{}, fmt.Errorf("error scan user")
	}
	return db.ConvertUserRepoToUserServ(user), nil
}

func (r *repo) GetById(ctx context.Context, id int64) (models.User, error) {
	builder := squirrel.Select(idColumn, loginColumn, passwordHashColumn, createdAtColumn).
		From(tableName).
		Where(squirrel.Eq{idColumn: id}).
		Limit(1).
		PlaceholderFormat(squirrel.Dollar)
	query, args, err := builder.ToSql()
	if err != nil {
		return models.User{}, err
	}
	row := r.db.QueryRow(ctx, query, args...)
	var user modelsRepo.User
	err = row.Scan(user.Id, user.Login, user.Password, user.CreatedAt)
	if err != nil {
		return models.User{}, fmt.Errorf("error scan user")
	}
	return db.ConvertUserRepoToUserServ(user), nil
}
