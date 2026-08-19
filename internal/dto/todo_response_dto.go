package dto

type ListResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ListsResponse struct {
	Data []ListResponse `json:"data"`
}

type ItemResponse struct {
	ID          int    `json:"id"`
	ListID      int    `json:"list_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type CreatedListResponse struct {
	ID int `json:"id"`
}

type CreatedItemResponse struct {
	ID int `json:"id"`
}
