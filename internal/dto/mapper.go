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

func (r CreateItemRequest) ToDomain(listID int) domain.TodoItem {
	return domain.TodoItem{
		ListID:      listID,
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

func NewItemResponse(i domain.TodoItem) ItemResponse {
	return ItemResponse{
		ID:          i.ID,
		ListID:      i.ListID,
		Title:       i.Title,
		Description: i.Description,
		Done:        i.Done,
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
