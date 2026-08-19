package dto

import (
	"github.com/IBKnight/todo-backend/internal/domain"
)

func (r CreateListRequest) ToDomain() domain.TodoList {
	return domain.TodoList{
		Title:       r.Title,
		Description: r.Description,
	}
}

func NewListResponse(l domain.TodoList) ListResponse {
	return ListResponse{
		ID:          l.ID,
		Title:       l.Title,
		Description: l.Description,
	}
}

func NewCreatedListsResponse(list domain.TodoList) CreatedListResponse {
	return CreatedListResponse{ID: list.ID}
}

func NewListsResponse(lists []domain.TodoList) ListsResponse {
	responseLists := make([]ListResponse, 0, len(lists))

	for i := range lists {
		responseLists = append(responseLists, ListResponse{
			ID:          lists[i].ID,
			Title:       lists[i].Title,
			Description: lists[i].Description,
		})
	}

	return ListsResponse{Data: responseLists}
}
