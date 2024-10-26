package adapter

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/protomem/chatik/internal/core/data"
	"github.com/protomem/chatik/internal/core/entity"
	"github.com/protomem/chatik/internal/infra/database"
	"github.com/protomem/chatik/pkg/werrors"
)

var (
	_ data.UserAccessor = (*UserDAO)(nil)
	_ data.UserMutator  = (*UserDAO)(nil)
)

type UserDAO struct {
	db      *database.DB
	builder database.QueryBuilder
}

func NewUserDAO(db *database.DB) *UserDAO {
	return &UserDAO{
		db:      db,
		builder: db.QueryBuilder(),
	}
}

func (dao *UserDAO) Select(ctx context.Context) ([]entity.User, error) {
	werr := werrors.Wrap("dao/user", "select")

	query, args, err := dao.
		querySelect().
		OrderBy("updated_at DESC").
		ToSql()
	if err != nil {
		return nil, werr(err, "build query")
	}

	rows, err := dao.db.Query(ctx, query, args...)
	if err != nil {
		return nil, werr(err)
	}
	defer rows.Close()

	users := make([]entity.User, 0)
	for rows.Next() {
		user, err := dao.scanUser(rows)
		if err != nil {
			return nil, werr(err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (dao *UserDAO) Get(ctx context.Context, id entity.ID) (entity.User, error) {
	werr := werrors.Wrap("dao/user", "get")

	query, args, err := dao.
		querySelect().
		Where(squirrel.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return entity.User{}, werr(err, "build query")
	}

	row := dao.db.QueryRow(ctx, query, args...)
	user, err := dao.scanUser(row)
	if err != nil {
		if database.IsNoRows(err) {
			return entity.User{}, werr(entity.ErrUserNotFound)
		}

		return entity.User{}, werr(err)
	}

	return user, nil
}

func (dao *UserDAO) GetByNickname(ctx context.Context, nickname string) (entity.User, error) {
	werr := werrors.Wrap("dao/user", "getByNickname")

	query, args, err := dao.
		querySelect().
		Where(squirrel.Eq{"nickname": nickname}).
		Limit(1).
		ToSql()
	if err != nil {
		return entity.User{}, werr(err, "build query")
	}

	row := dao.db.QueryRow(ctx, query, args...)
	user, err := dao.scanUser(row)
	if err != nil {
		if database.IsNoRows(err) {
			return entity.User{}, werr(entity.ErrUserNotFound)
		}

		return entity.User{}, werr(err)
	}

	return user, nil
}

func (dao *UserDAO) Insert(ctx context.Context, dto data.InsertUserDTO) (entity.ID, error) {
	werr := werrors.Wrap("dao/user", "insert")

	query, args, err := dao.
		queryInsert(dto).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return entity.ID{}, werr(err, "build query")
	}

	var id entity.ID
	row := dao.db.QueryRow(ctx, query, args...)
	if err := row.Scan(&id); err != nil {
		if dberr, ok := database.AsUniqueConstraint(err); ok {
			return entity.ID{}, werr(entity.ErrUserExists.
				WithField(ExtractColumnFromContraint(dberr.ConstraintName)))
		}

		return entity.ID{}, werr(err)
	}

	return id, nil
}

func (dao *UserDAO) querySelect() squirrel.SelectBuilder {
	return dao.builder.
		Select(
			"id", "created_at", "updated_at",
			"nickname", "password",
			"email", "is_verified",
		).
		From("users")
}

func (dao *UserDAO) queryInsert(dto data.InsertUserDTO) squirrel.InsertBuilder {
	return dao.builder.
		Insert("users").
		Columns("nickname", "password", "email").
		Values(dto.Nickname, dto.Password, dto.Email)
}

func (*UserDAO) scanUser(s database.Scanner) (entity.User, error) {
	var user entity.User
	return user, s.Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
		&user.Nickname, &user.Password,
		&user.Email, &user.Verified,
	)
}
