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
	_ data.MessageAccessor = (*MessageDAO)(nil)
	_ data.MessageMutator  = (*MessageDAO)(nil)
)

type MessageDAO struct {
	db      *database.DB
	builder database.QueryBuilder
}

func NewMessageDAO(db *database.DB) *MessageDAO {
	return &MessageDAO{
		db:      db,
		builder: db.QueryBuilder(),
	}
}

func (dao *MessageDAO) SelectByFromAndTo(ctx context.Context, opts data.SelectMessageByFromAndToOptions) ([]entity.Message, error) {
	werr := werrors.Wrap("dao/messages", "findByFromAndTo")

	query, args, err := dao.
		querySelect().
		Where(squirrel.Eq{"messages.from_user_id": opts.From, "messages.to_user_id": opts.To}).
		OrderBy("messages.created_at DESC").
		Limit(uint64(opts.Limit)).
		Offset(uint64(opts.Offset)).
		ToSql()
	if err != nil {
		return []entity.Message{}, werr(err)
	}

	rows, err := dao.db.Query(ctx, query, args...)
	if err != nil {
		return []entity.Message{}, werr(err)
	}
	defer rows.Close()

	messages := make([]entity.Message, 0, opts.Limit)
	for rows.Next() {
		message, err := dao.scanMessage(rows)
		if err != nil {
			return []entity.Message{}, werr(err)
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (dao *MessageDAO) Get(ctx context.Context, id entity.ID) (entity.Message, error) {
	werr := werrors.Wrap("dao/message", "get")

	query, args, err := dao.
		querySelect().
		Where(squirrel.Eq{"messages.id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return entity.Message{}, werr(err, "build query")
	}

	row := dao.db.QueryRow(ctx, query, args...)
	message, err := dao.scanMessage(row)
	if err != nil {
		if database.IsNoRows(err) {
			return entity.Message{}, werr(entity.ErrMessageNotFound)
		}

		return entity.Message{}, werr(err)
	}

	return message, nil
}

func (dao *MessageDAO) Insert(ctx context.Context, dto data.InsertMessageDTO) (entity.ID, error) {
	werr := werrors.Wrap("dao/message", "insert")

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
		return entity.ID{}, werr(err)
	}

	return id, nil
}

func (dao *MessageDAO) querySelect() squirrel.SelectBuilder {
	return dao.builder.
		Select(
			"messages.id", "messages.created_at", "messages.updated_at",
			"messages.message_text", "messages.metadata",

			"from_users.id", "from_users.created_at", "from_users.updated_at",
			"from_users.nickname", "from_users.password", "from_users.email", "from_users.is_verified",

			"to_users.id", "to_users.created_at", "to_users.updated_at",
			"to_users.nickname", "to_users.password", "to_users.email", "to_users.is_verified",
		).
		From("messages").
		Join("users as from_users ON from_users.id = messages.from_user_id").
		Join("users as to_users ON to_users.id = messages.to_user_id")
}

func (dao *MessageDAO) queryInsert(dto data.InsertMessageDTO) squirrel.InsertBuilder {
	return dao.builder.
		Insert("messages").
		Columns("from_user_id", "to_user_id", "message_text", "metadata").
		Values(dto.From, dto.To, dto.Text, dto.Metadata)
}

func (*MessageDAO) scanMessage(s database.Scanner) (entity.Message, error) {
	var message entity.Message
	return message, s.Scan(
		&message.ID, &message.CreatedAt, &message.UpdatedAt,
		&message.Text, &message.Metadata,

		&message.From.ID, &message.From.CreatedAt, &message.From.UpdatedAt,
		&message.From.Nickname, &message.From.Password, &message.From.Email, &message.From.Verified,

		&message.To.ID, &message.To.CreatedAt, &message.To.UpdatedAt,
		&message.To.Nickname, &message.To.Password, &message.To.Email, &message.To.Verified,
	)
}
