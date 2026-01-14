package models

type Task struct {
	ID    int32  `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
}
