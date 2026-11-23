package model

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	IsDone bool   `json:"isDone"`
}

type TaskPatch struct {
	Title  *string `json:"title,omitempty"`
	IsDone *bool   `json:"isDone,omitempty"`
}
