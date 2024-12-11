package document

import (
	"context"
	"doc-server/internal/models"
	"doc-server/internal/storage"
	"doc-server/internal/storage/db"
	models2 "doc-server/internal/storage/db/pg/models"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	tableNameDocument              = "documents"
	tableNameGrant                 = "grants"
	idDocumentColumn               = "id"
	nameDocumentColumn             = "name"
	mimeDocumentColumn             = "mime"
	isPublicDocumentDocumentColumn = "is_public"
	isFileDocumentDocumentColumn   = "is_file"
	userIdDocumentColumn           = "user_id"
	createdAtDocumentColumn        = "created_at"
	idGrantColumn                  = "id"
	loginGrantColumn               = "login"
	documentIdGrantColumn          = "document_id"
	permissionGrantColumn          = "permission"
	createdAtGrantColumn           = "created_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewDocumentRepository(db *pgxpool.Pool) storage.DocumentRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, document models.Document) error {
	builder := squirrel.Insert(tableNameDocument).
		Columns(
			nameDocumentColumn,
			mimeDocumentColumn,
			isPublicDocumentDocumentColumn,
			isFileDocumentDocumentColumn,
			userIdDocumentColumn,
		).
		Values(
			document.Name,
			document.Mime,
			document.IsPublic,
			document.IsFile,
			document.UserId,
		).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("returning id")

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return err
	}
	var document_id int64
	row := tx.QueryRow(ctx, query, args...)
	err = row.Scan(&document_id)
	if err != nil {
		return err
	}
	for i := 0; i < len(document.Grant); i++ {
		builder = squirrel.Insert(tableNameGrant).
			Columns(
				loginGrantColumn,
				documentIdGrantColumn,
				permissionGrantColumn,
			).
			Values(document.Grant[i].Login,
				document_id,
				document.Grant[i].Permission).
			PlaceholderFormat(squirrel.Dollar)
		query, args, err = builder.ToSql()
		tag, err := tx.Exec(ctx, query, args...)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return fmt.Errorf("error: failed save grant")
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) GetList(ctx context.Context, userId int64, filter map[string]interface{}, limit int64) ([]models.Document, error) {
	// Билдер для основного запроса
	builder := squirrel.Select(
		fmt.Sprintf("d.%s AS document_id", idDocumentColumn),
		fmt.Sprintf("d.%s", nameDocumentColumn),
		fmt.Sprintf("d.%s", mimeDocumentColumn),
		fmt.Sprintf("d.%s", isPublicDocumentDocumentColumn),
		fmt.Sprintf("d.%s", isFileDocumentDocumentColumn),
		fmt.Sprintf("d.%s", userIdDocumentColumn),
		fmt.Sprintf("d.%s AS document_created_at", createdAtDocumentColumn),
		fmt.Sprintf("g.%s AS grant_id", idGrantColumn),
		fmt.Sprintf("g.%s", loginGrantColumn),
		fmt.Sprintf("g.%s", documentIdGrantColumn),
		fmt.Sprintf("g.%s", permissionGrantColumn),
		fmt.Sprintf("g.%s AS grant_created_at", createdAtGrantColumn),
	).
		From(fmt.Sprintf("%s d", tableNameDocument)).
		Join(fmt.Sprintf("%s g ON d.%s = g.%s", tableNameGrant, idDocumentColumn, documentIdGrantColumn)).
		Where(squirrel.Eq{fmt.Sprintf("d.%s", userIdDocumentColumn): userId}).
		Limit(uint64(limit)).
		PlaceholderFormat(squirrel.Dollar)

	// Применяем фильтры из мапы
	for key, value := range filter {
		switch key {
		case "login":
			builder = builder.Where(squirrel.Eq{fmt.Sprintf("g.%s", loginGrantColumn): value})
		case "mime":
			builder = builder.Where(squirrel.Eq{fmt.Sprintf("d.%s", mimeDocumentColumn): value})
		case "is_public":
			builder = builder.Where(squirrel.Eq{fmt.Sprintf("d.%s", isPublicDocumentDocumentColumn): value})
		case "is_file":
			builder = builder.Where(squirrel.Eq{fmt.Sprintf("d.%s", isFileDocumentDocumentColumn): value})
		}
	}

	builder = builder.OrderBy(
		fmt.Sprintf("d.%s ASC", nameDocumentColumn),      // по имени в порядке возрастания
		fmt.Sprintf("d.%s ASC", createdAtDocumentColumn), // по дате создания в порядке возрастания
	)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	documentMap := make(map[int64]*models2.Document)

	for rows.Next() {
		var document models2.Document
		var grant models2.Grant

		err = rows.Scan(
			&document.Id,
			&document.Name,
			&document.Mime,
			&document.IsPublic,
			&document.IsFile,
			&document.UserId,
			&document.CreatedAt,
			&grant.Id,
			&grant.Login,
			&grant.DocumentId,
			&grant.Permission,
			&grant.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if existingDocument, exists := documentMap[document.Id]; exists {
			existingDocument.Grant = append(existingDocument.Grant, grant)
		} else {
			document.Grant = []models2.Grant{grant}
			documentMap[document.Id] = &document
		}
	}

	documents := []models.Document{}
	for _, doc := range documentMap {
		documents = append(documents, db.ConverterDocumentRepoToDocumentServ(*doc))
	}

	return documents, nil
}

func (r *repo) GetPublicList(ctx context.Context, filter map[string]interface{}, limit int64) ([]models.Document, error) {
	builder := squirrel.Select(
		fmt.Sprintf("d.%s AS document_id", idDocumentColumn),
		fmt.Sprintf("d.%s", nameDocumentColumn),
		fmt.Sprintf("d.%s", mimeDocumentColumn),
		fmt.Sprintf("d.%s", isPublicDocumentDocumentColumn),
		fmt.Sprintf("d.%s", isFileDocumentDocumentColumn),
		fmt.Sprintf("d.%s", userIdDocumentColumn),
		fmt.Sprintf("d.%s AS document_created_at", createdAtDocumentColumn),
		fmt.Sprintf("g.%s AS grant_id", idGrantColumn),
		fmt.Sprintf("g.%s", loginGrantColumn),
		fmt.Sprintf("g.%s", documentIdGrantColumn),
		fmt.Sprintf("g.%s", permissionGrantColumn),
		fmt.Sprintf("g.%s AS grant_created_at", createdAtGrantColumn),
	).
		From(fmt.Sprintf("%s d", tableNameDocument)).
		LeftJoin(fmt.Sprintf("%s g ON d.%s = g.%s", tableNameGrant, idDocumentColumn, documentIdGrantColumn)).
		Where(squirrel.Eq{fmt.Sprintf("d.%s", isPublicDocumentDocumentColumn): true}).
		Limit(uint64(limit)).
		PlaceholderFormat(squirrel.Dollar)

	for key, value := range filter {
		switch key {
		case "mime":
			builder = builder.Where(squirrel.Eq{fmt.Sprintf("d.%s", mimeDocumentColumn): value})
		case "is_file":
			builder = builder.Where(squirrel.Eq{fmt.Sprintf("d.%s", isFileDocumentDocumentColumn): value})
		}
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Используем мапу для сопоставления документов с их грантами
	documentMap := make(map[int64]*models2.Document)

	for rows.Next() {
		var document models2.Document
		var grant models2.Grant

		err = rows.Scan(
			&document.Id,
			&document.Name,
			&document.Mime,
			&document.IsPublic,
			&document.IsFile,
			&document.UserId,
			&document.CreatedAt,
			&grant.Id,
			&grant.Login,
			&grant.DocumentId,
			&grant.Permission,
			&grant.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Если документ уже есть в мапе, добавляем grant
		if existingDocument, exists := documentMap[document.Id]; exists {
			existingDocument.Grant = append(existingDocument.Grant, grant)
		} else {
			// Если документа нет, создаем новый
			document.Grant = []models2.Grant{grant}
			documentMap[document.Id] = &document
		}
	}

	// Преобразуем мапу в слайс
	documents := []models.Document{}
	for _, doc := range documentMap {
		documents = append(documents, db.ConverterDocumentRepoToDocumentServ(*doc))
	}

	return documents, nil
}

func (r *repo) Get(ctx context.Context, login string, documentId int64) (models.Document, error) {
	builder := squirrel.Select(
		fmt.Sprintf("d.%s AS document_id", idDocumentColumn),
		fmt.Sprintf("d.%s", nameDocumentColumn),
		fmt.Sprintf("d.%s", mimeDocumentColumn),
		fmt.Sprintf("d.%s", isPublicDocumentDocumentColumn),
		fmt.Sprintf("d.%s", userIdDocumentColumn),
		fmt.Sprintf("d.%s AS document_created_at", createdAtDocumentColumn),
		fmt.Sprintf("g.%s AS grant_id", idGrantColumn),
		fmt.Sprintf("g.%s", loginGrantColumn),
		fmt.Sprintf("g.%s", documentIdGrantColumn),
		fmt.Sprintf("g.%s", permissionGrantColumn),
		fmt.Sprintf("g.%s AS grant_created_at", createdAtGrantColumn),
	).
		From(fmt.Sprintf("%s d", tableNameDocument)).
		LeftJoin(fmt.Sprintf("%s g ON d.%s = g.%s", tableNameGrant, idDocumentColumn, documentIdGrantColumn)).
		Where(squirrel.Eq{fmt.Sprintf("d.%s", idDocumentColumn): documentId}).
		Where(squirrel.Or{
			squirrel.Eq{fmt.Sprintf("d.%s", isPublicDocumentDocumentColumn): true},
			squirrel.Eq{fmt.Sprintf("g.%s", loginGrantColumn): login},
		}).
		Limit(1).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return models.Document{}, fmt.Errorf("error building query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return models.Document{}, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return models.Document{}, fmt.Errorf("document not found or access denied")
	}

	var document models2.Document
	var grant models2.Grant

	err = rows.Scan(
		&document.Id,
		&document.Name,
		&document.Mime,
		&document.IsPublic,
		&document.UserId,
		&document.CreatedAt,
		&grant.Id,
		&grant.Login,
		&grant.DocumentId,
		&grant.Permission,
		&grant.CreatedAt)
	if err != nil {
		return models.Document{}, fmt.Errorf("error scanning rows: %w", err)
	}

	document.Grant = append(document.Grant, grant)
	return db.ConverterDocumentRepoToDocumentServ(document), nil
}

func (r *repo) Delete(ctx context.Context, userId int64, documentId int64) (string, error) {
	builder := squirrel.Delete(tableNameDocument).
		Where(squirrel.And{
			squirrel.Eq{userIdDocumentColumn: userId},
			squirrel.Eq{idDocumentColumn: documentId},
		}).Suffix("RETURNING name").
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return "", err
	}

	var name string
	err = r.db.QueryRow(ctx, query, args...).Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}
