package database

import "github.com/alimgiray/tasks/src/models"

func (db *Database) CreateTask(task *models.Task) error {
	query := `
		INSERT INTO tasks (key, deadline, assignee, title, status, description, project_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.db.Exec(query, task.Key, task.Deadline, task.Assignee, task.Title, task.Status, task.Description, task.ProjectID)
	if err != nil {
		return err
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	task.TaskID = int(taskID)
	return nil
}

func (db *Database) UpdateTask(task *models.Task) error {
	query := `
		UPDATE tasks
		SET key = ?, deadline = ?, assignee = ?, title = ?, status = ?, description = ?, project_id = ?
		WHERE id = ?
	`
	_, err := db.db.Exec(query, task.Key, task.Deadline, task.Assignee, task.Title, task.Status, task.Description, task.ProjectID, task.TaskID)
	return err
}

func (db *Database) GetTasks(projectID int) ([]*models.Task, error) {
	query := `
		SELECT key, deadline, created_at, assignee, title, status, description, project_id
		FROM tasks
		WHERE project_id = ?
	`

	rows, err := db.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.Key, &task.Deadline, &task.CreatedAt, &task.Assignee, &task.Title, &task.Status, &task.Description, &task.ProjectID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func LoadTasksFromDB(projectID int) []*models.Task {
	tasks, err := Get().GetTasks(projectID)
	if err != nil {
		panic(err)
	}
	return tasks
}
