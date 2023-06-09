package model

// Todo — модель задачи.
type Todo struct {
	ID       uint32 `json:"id"`
	Name     string `json:"name"`
	Done     bool   `json:"done"`
	Priority uint8  `json:"priority"`
}
