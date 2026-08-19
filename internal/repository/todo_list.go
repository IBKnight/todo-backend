package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/IBKnight/todo-backend/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TodoListRepo struct {
	db *sqlx.DB
}

func NewTodoListRepo(db *sqlx.DB) *TodoListRepo {
	return &TodoListRepo{db}
}

func (r *TodoListRepo) CreateList(ctx context.Context, userId int, list domain.TodoList) (int, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var listId int
	createListQuery := fmt.Sprintf(
		"INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id",
		todoListsTable,
	)
	if err := tx.QueryRowContext(ctx, createListQuery, list.Title, list.Description).Scan(&listId); err != nil {
		return 0, fmt.Errorf("insert list: %w", err)
	}

	createUsersListQuery := fmt.Sprintf(
		"INSERT INTO %s (user_id, list_id) VALUES ($1, $2)",
		usersListsTable,
	)
	if _, err := tx.ExecContext(ctx, createUsersListQuery, userId, listId); err != nil {
		return 0, fmt.Errorf("insert users_list: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit tx: %w", err)
	}

	return listId, nil
}

func (r *TodoListRepo) GetUserLists(ctx context.Context, userId int) ([]domain.TodoList, error) {
	query := fmt.Sprintf(
		`SELECT tl.id, tl.title, tl.description
		 FROM %s tl
		 INNER JOIN %s ul ON tl.id = ul.list_id
		 WHERE ul.user_id = $1
		 ORDER BY tl.id`,
		todoListsTable, usersListsTable,
	)

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("query lists: %w", err)
	}
	defer rows.Close()

	lists := make([]domain.TodoList, 0)
	for rows.Next() {
		var l domain.TodoList
		if err := rows.Scan(&l.ID, &l.Title, &l.Description); err != nil {
			return nil, fmt.Errorf("scan list: %w", err)
		}
		lists = append(lists, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return lists, nil
}

func (r *TodoListRepo) GetListById(ctx context.Context, listId int, userId int) (domain.TodoList, error) {
	var list todoListRow

	query := fmt.Sprintf(
		`SELECT tl.id, tl.title, tl.description
		 FROM %s tl
		 INNER JOIN %s ul ON tl.id = ul.list_id
		 WHERE ul.user_id = $1 AND ul.list_id = $2`,
		todoListsTable, usersListsTable,
	)

	row := r.db.QueryRowContext(ctx, query, userId, listId)
	err := row.Scan(&list.ID, &list.Title, &list.Description)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.TodoList{}, domain.ErrListNotFound
	}
	if err != nil {
		return domain.TodoList{}, fmt.Errorf("scan list: %w", err)
	}

	return list.toDomain(), nil
}

func (r *TodoListRepo) RemoveList(ctx context.Context, userId int, listId int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s tl
		 USING %s ul
		 WHERE tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2`,
		todoListsTable, usersListsTable,
	)

	res, err := r.db.ExecContext(ctx, query, userId, listId)
	if err != nil {
		return fmt.Errorf("delete list: %w", err)
	}

	logrus.Info(res)

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if affected == 0 {
		return domain.ErrListNotFound
	}

	return nil
}

func (r *TodoListRepo) UpdateList(ctx context.Context, userId int, updatedTodoList domain.TodoList) (domain.TodoList, error) {
	row := toRow(updatedTodoList)

	query := fmt.Sprintf(
		`UPDATE %s tl SET title=$1, description=$2
	 FROM %s ul
	 WHERE tl.id = ul.list_id AND ul.list_id=$3 AND ul.user_id=$4
	 RETURNING tl.id, tl.title, tl.description`,
		todoListsTable, usersListsTable,
	)

	var updated domain.TodoList

	err := r.db.QueryRowContext(ctx, query,
		row.Title, row.Description, row.ID, userId,
	).Scan(&updated.ID, &updated.Title, &updated.Description)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.TodoList{}, domain.ErrListNotFound
	}
	if err != nil {
		return domain.TodoList{}, fmt.Errorf("update list: %w", err)
	}

	return updated, nil
}
