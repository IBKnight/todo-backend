package dto

import "github.com/IBKnight/todo-backend/internal/domain"

func (r CreateItemRequest) ToDomain() domain.TodoItem {
	return domain.TodoItem{
		Title:       r.Title,
		Description: r.Description,
	}
}

func (r UpdateItemRequest) ToDomain(id int) domain.TodoItem {
	return domain.TodoItem{
		ID:          id,
		Title:       r.Title,
		Description: r.Description,
		Done:        r.Done,
	}
}

func NewItemResponse(i domain.TodoItem) ItemResponse {
	return ItemResponse{
		ID:          i.ID,
		Title:       i.Title,
		Description: i.Description,
		Done:        i.Done,
	}
}

func NewItemsResponse(items []domain.TodoItem) ItemsResponse {
	out := make([]ItemResponse, 0, len(items))
	for _, i := range items {
		out = append(out, NewItemResponse(i))
	}
	return ItemsResponse{Data: out}
}
