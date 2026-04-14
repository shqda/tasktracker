package postgres

import (
	"github.com/jmoiron/sqlx"
	"tasktracker/internal/model"
)

type PostgresDB struct {
	DB *sqlx.DB
}

func (p *PostgresDB) InsertTask(title string) (int, error) {
	var id int
	if err := p.DB.Get(&id, "insert into tasks (title) values ($1) returning id", title); err != nil {
		return 0, err
	}
	return id, nil
}

func (p *PostgresDB) GetLastTask() (*model.Task, error) {
	var task model.Task
	if err := p.DB.Get(&task, "select id, title from tasks order by id desc limit 1"); err != nil {
		return nil, err
	}
	return &task, nil
}
