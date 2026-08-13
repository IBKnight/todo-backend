package domain

type TodoList struct {
	Id          int    `json:"-"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoItem struct {
	Id          int    `json:"-"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}
type UsersList struct {
	Id     int `json:"-"`
	UserId int `json:"user_id"`
	ListId int `json:"list_id"`
}
type ListItems struct {
	Id     int `json:"-"`
	ListId int `json:"list_id"`
	ItemId int `json:"item_id"`
}
