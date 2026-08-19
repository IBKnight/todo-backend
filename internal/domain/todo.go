package domain

type TodoList struct {
	ID          int
	Title       string
	Description string
}

type TodoItem struct {
	ID          int
	ListID      int
	Title       string
	Description string
	Done        bool
}

type UsersList struct {
	ID     int
	UserID int
	ListID int
}

type ListItem struct {
	ID     int
	ListID int
	ItemID int
}
