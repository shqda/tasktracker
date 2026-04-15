package postgres

import (
	"tasktracker/internal/model"

	"github.com/jmoiron/sqlx"
)

type PostgresDB struct {
	DB *sqlx.DB
}

func (p *PostgresDB) GetLastTask() (*model.Task, error) {
	var task model.Task
	if err := p.DB.Get(&task, "select id, title from tasks order by id desc limit 1"); err != nil {
		return nil, err
	}
	return &task, nil
}

func (p *PostgresDB) GetTaskByID(id int) (*model.Task, error) {
	var task model.Task
	if err := p.DB.Get(&task, "select id, title from tasks where id = $1", id); err != nil {
		return nil, err
	}
	return &task, nil
}

func (p *PostgresDB) GetAllTasks() ([]model.Task, error) {
	var tasks []model.Task
	if err := p.DB.Select(&tasks, "select id, title from tasks"); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (p *PostgresDB) DeleteTask(id int) error {
	if _, err := p.DB.Exec("delete from tasks where id = $1", id); err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) UpdateTask(id int, title string) error {
	if _, err := p.DB.Exec("update tasks set title = $1 where id = $2", title, id); err != nil {
		return err
	}
	return nil
}
func (p *PostgresDB) InsertTask(title string) (int, error) {
	var id int
	if err := p.DB.Get(&id, "insert into tasks (title) values ($1) returning id", title); err != nil {
		return 0, err
	}
	return id, nil
}
