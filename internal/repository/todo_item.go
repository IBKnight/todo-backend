package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/IBKnight/todo-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type TodoItemRepo struct {
	db *sqlx.DB
}

func NewTodoItemRepo(db *sqlx.DB) *TodoItemRepo {
	return &TodoItemRepo{db}
}

func (r *TodoItemRepo) Create(ctx context.Context, userID, listID int, item domain.TodoItem) (domain.TodoItem, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return domain.TodoItem{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var created domain.TodoItem
	createItemQuery := fmt.Sprintf(
		`INSERT INTO %s (title, description, done) VALUES ($1, $2, $3)
		 RETURNING id, title, description, done`,
		todoItemsTable,
	)
	err = tx.QueryRowContext(ctx, createItemQuery, item.Title, item.Description, item.Done).
		Scan(&created.ID, &created.Title, &created.Description, &created.Done)
	if err != nil {
		return domain.TodoItem{}, fmt.Errorf("insert item: %w", err)
	}

	linkQuery := fmt.Sprintf(
		`INSERT INTO %s (list_id, item_id)
		 SELECT $1, $2
		 WHERE EXISTS (SELECT 1 FROM %s WHERE list_id = $1 AND user_id = $3)`,
		listItemsTable, usersListsTable,
	)
	res, err := tx.ExecContext(ctx, linkQuery, listID, created.ID, userID)
	if err != nil {
		return domain.TodoItem{}, fmt.Errorf("insert list_item: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return domain.TodoItem{}, fmt.Errorf("rows affected: %w", err)
	}
	if affected == 0 {
		return domain.TodoItem{}, domain.ErrListNotFound
	}

	if err := tx.Commit(); err != nil {
		return domain.TodoItem{}, fmt.Errorf("commit tx: %w", err)
	}

	created.ListID = listID
	return created, nil
}

func (r *TodoItemRepo) GetAll(ctx context.Context, userID, listID int) ([]domain.TodoItem, error) {
	query := fmt.Sprintf(
		`SELECT ti.id, li.list_id, ti.title, ti.description, ti.done
		 FROM %s ti
		 INNER JOIN %s li ON ti.id = li.item_id
		 INNER JOIN %s ul ON li.list_id = ul.list_id
		 WHERE li.list_id = $1 AND ul.user_id = $2
		 ORDER BY ti.id`,
		todoItemsTable, listItemsTable, usersListsTable,
	)

	rows, err := r.db.QueryContext(ctx, query, listID, userID)
	if err != nil {
		return nil, fmt.Errorf("query items: %w", err)
	}
	defer rows.Close()

	items := make([]domain.TodoItem, 0)
	for rows.Next() {
		var i domain.TodoItem
		if err := rows.Scan(&i.ID, &i.ListID, &i.Title, &i.Description, &i.Done); err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}
		items = append(items, i)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *TodoItemRepo) GetByID(ctx context.Context, userID, itemID int) (domain.TodoItem, error) {
	query := fmt.Sprintf(
		`SELECT ti.id, li.list_id, ti.title, ti.description, ti.done
		 FROM %s ti
		 INNER JOIN %s li ON ti.id = li.item_id
		 INNER JOIN %s ul ON li.list_id = ul.list_id
		 WHERE ti.id = $1 AND ul.user_id = $2`,
		todoItemsTable, listItemsTable, usersListsTable,
	)

	var item domain.TodoItem
	err := r.db.QueryRowContext(ctx, query, itemID, userID).
		Scan(&item.ID, &item.ListID, &item.Title, &item.Description, &item.Done)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.TodoItem{}, domain.ErrItemNotFound
	}
	if err != nil {
		return domain.TodoItem{}, fmt.Errorf("get item: %w", err)
	}

	return item, nil
}

func (r *TodoItemRepo) Update(ctx context.Context, userID int, item domain.TodoItem) (domain.TodoItem, error) {
	query := fmt.Sprintf(
		`UPDATE %s ti SET title=$1, description=$2, done=$3
		 FROM %s li, %s ul
		 WHERE ti.id = li.item_id AND li.list_id = ul.list_id
		   AND ti.id = $4 AND ul.user_id = $5
		 RETURNING ti.id, li.list_id, ti.title, ti.description, ti.done`,
		todoItemsTable, listItemsTable, usersListsTable,
	)

	var updated domain.TodoItem
	err := r.db.QueryRowContext(ctx, query,
		item.Title, item.Description, item.Done, item.ID, userID,
	).Scan(&updated.ID, &updated.ListID, &updated.Title, &updated.Description, &updated.Done)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.TodoItem{}, domain.ErrItemNotFound
	}
	if err != nil {
		return domain.TodoItem{}, fmt.Errorf("update item: %w", err)
	}

	return updated, nil
}

func (r *TodoItemRepo) Delete(ctx context.Context, userID, itemID int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s ti
		 USING %s li, %s ul
		 WHERE ti.id = li.item_id AND li.list_id = ul.list_id
		   AND ti.id = $1 AND ul.user_id = $2`,
		todoItemsTable, listItemsTable, usersListsTable,
	)

	res, err := r.db.ExecContext(ctx, query, itemID, userID)
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrItemNotFound
	}

	return nil
}
