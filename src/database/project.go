package database

import "github.com/alimgiray/tasks/src/models"

func (db *Database) CreateProject(project *models.Project) (int, error) {
	query := `
		INSERT INTO projects (name, created_at, updated_at)
		VALUES (?, ?, ?)
	`
	result, err := db.db.Exec(query, project.Name, project.CreatedAt, project.UpdatedAt)
	if err != nil {
		return 0, err
	}

	projectID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	project.ProjectID = int(projectID)
	return project.ProjectID, nil
}

func (db *Database) UpdateProject(project *models.Project) error {
	query := `
		UPDATE projects
		SET name = ?
		WHERE id = ?
	`
	_, err := db.db.Exec(query, project.Name, project.ProjectID)
	return err
}

func (db *Database) GetProjects() ([]*models.Project, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM projects
	`

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var project models.Project
		err := rows.Scan(&project.ProjectID, &project.Name, &project.CreatedAt, &project.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &project)
	}

	return projects, nil
}

func LoadProjectsFromDB() []*models.Project {
	projects, err := Get().GetProjects()
	if err != nil {
		panic(err)
	}
	return projects
}
