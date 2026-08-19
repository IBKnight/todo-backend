package dto

type CreateListRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Description string `json:"description" binding:"max=1000"`
}

type UpdateListRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Description string `json:"description"`
}

type ListResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ListsResponse struct {
	Data []ListResponse `json:"data"`
}

type CreatedListResponse struct {
	ID int `json:"id"`
}
