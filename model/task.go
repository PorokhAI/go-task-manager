package model

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	IsDone      bool   `json:"isDone"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
}

type TaskPatch struct {
	Title       *string `json:"title,omitempty"`
	IsDone      *bool   `json:"isDone,omitempty"`
	Description *string `json:"description,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
}
