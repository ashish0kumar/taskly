package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	gap "github.com/muesli/go-app-paths"

	"github.com/ashish0kumar/taskly/internal/task"
)

// TaskDB holds the database connection. Exported type.
type TaskDB struct {
	db *sql.DB
}

// setupPath determines the application's data directory path. (Unexported)
func setupPath() (string, error) {
	scope := gap.NewScope(gap.User, "taskly")
	dirs, err := scope.DataDirs()
	if err != nil {
		return "", fmt.Errorf("could not determine data directory: %w", err)
	}
	var taskDir string
	if len(dirs) > 0 {
		taskDir = dirs[0]
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not determine home directory as fallback: %w", err)
		}
		taskDir = filepath.Join(home, ".taskly")
	}
	if err := initTaskDir(taskDir); err != nil {
		return "", fmt.Errorf("could not initialize task directory '%s': %w", taskDir, err)
	}
	return taskDir, nil
}

// initTaskDir ensures the application's data directory exists. (Unexported)
func initTaskDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, 0o750)
		}
		return err
	}
	return nil
}

// GetStoragePath returns the calculated storage path for the database file. Exported
func GetStoragePath() (string, error) {
	dataDir, err := setupPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, "tasks.db"), nil
}

// OpenDB establishes a connection and ensures the table exists. Exported
func OpenDB() (*TaskDB, error) {
	// Use GetStoragePath to determine the full db file path
	dbPath, err := GetStoragePath()
	if err != nil {
		return nil, fmt.Errorf("could not get storage path: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database '%s': %w", dbPath, err)
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database '%s': %w", dbPath, err)
	}

	t := &TaskDB{db: db} // Create exported TaskDB

	exists, err := t.tableExists("tasks")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to check if 'tasks' table exists: %w", err)
	}
	if !exists {
		err := t.createTable()
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to create 'tasks' table: %w", err)
		}
	}
	return t, nil
}

// Close closes the database connection. Exported
func (tdb *TaskDB) Close() error {
	if tdb.db != nil {
		return tdb.db.Close()
	}
	return nil
}

// tableExists checks if the 'tasks' table exists. (Unexported)
func (tdb *TaskDB) tableExists(name string) (bool, error) {
	query := "SELECT name FROM sqlite_master WHERE type='table' AND name=?;"
	var tableName string
	err := tdb.db.QueryRow(query, name).Scan(&tableName)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed checking table existence: %w", err)
	}
	return tableName == name, nil
}

// createTable creates the 'tasks' table. (Unexported)
func (tdb *TaskDB) createTable() error {
	schema := `
	CREATE TABLE "tasks" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"name" TEXT NOT NULL CHECK(length(name) > 0),
		"project" TEXT,
		"status" TEXT NOT NULL DEFAULT 'todo' CHECK(status IN ('todo', 'in progress', 'done')),
		"created" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := tdb.db.Exec(schema)
	return err
}

// --- Exported CRUD Methods ---

// Insert adds a new task.
func (tdb *TaskDB) Insert(name, project string) (task.Task, error) {
	createdTime := time.Now()
	defaultStatus := task.Todo.String() // Use Status enum from task package

	stmt := "INSERT INTO tasks(name, project, status, created) VALUES(?, ?, ?, ?)"
	res, err := tdb.db.Exec(stmt, name, project, defaultStatus, createdTime)
	if err != nil {
		return task.Task{}, fmt.Errorf("insert failed: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return task.Task{}, fmt.Errorf("failed to get last insert ID: %w", err)
	}
	return task.Task{
		ID:      uint(id),
		Name:    name,
		Project: project,
		Status:  defaultStatus,
		Created: createdTime,
	}, nil
}

// Delete removes a task by ID.
func (tdb *TaskDB) Delete(id uint) error {
	res, err := tdb.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete failed for id %d: %w", id, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err == nil && rowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found for deletion", id)
	}
	return err
}

// Update modifies an existing task.
func (tdb *TaskDB) Update(id uint, name *string, project *string, status *string) (task.Task, error) {
	orig, err := tdb.GetTask(id) // Use exported GetTask
	if err != nil {
		return task.Task{}, fmt.Errorf("cannot update task %d: %w", id, err)
	}
	setClauses := []string{}
	args := []interface{}{}
	if name != nil {
		setClauses = append(setClauses, "name = ?")
		args = append(args, *name)
		orig.Name = *name
	}
	if project != nil {
		setClauses = append(setClauses, "project = ?")
		args = append(args, *project)
		orig.Project = *project
	}
	if status != nil {
		setClauses = append(setClauses, "status = ?")
		args = append(args, *status)
		orig.Status = *status
	}
	if len(setClauses) == 0 {
		return orig, nil
	}
	args = append(args, id)
	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = ?", strings.Join(setClauses, ", "))
	res, err := tdb.db.Exec(query, args...)
	if err != nil {
		return task.Task{}, fmt.Errorf("db update failed for id %d: %w", id, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err == nil && rowsAffected == 0 {
		return task.Task{}, fmt.Errorf("task with ID %d not found for update", id)
	}
	return orig, nil
}

// GetTasks retrieves all tasks.
func (tdb *TaskDB) GetTasks() ([]task.Task, error) {
	tasks := []task.Task{}
	rows, err := tdb.db.Query("SELECT id, name, project, status, created FROM tasks ORDER BY created ASC")
	if err != nil {
		return nil, fmt.Errorf("unable to query tasks: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var t task.Task
		var project sql.NullString
		err = rows.Scan(&t.ID, &t.Name, &project, &t.Status, &t.Created)
		if err != nil {
			return nil, fmt.Errorf("failed scanning task row: %w", err)
		}
		if project.Valid {
			t.Project = project.String
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating task rows: %w", err)
	}
	return tasks, nil
}

// GetTask retrieves a single task by ID.
func (tdb *TaskDB) GetTask(id uint) (task.Task, error) {
	var t task.Task
	var project sql.NullString
	err := tdb.db.QueryRow("SELECT id, name, project, status, created FROM tasks WHERE id = ?", id).
		Scan(&t.ID, &t.Name, &project, &t.Status, &t.Created)
	if err != nil {
		if err == sql.ErrNoRows {
			return task.Task{}, fmt.Errorf("task with ID %d not found", id)
		}
		return task.Task{}, fmt.Errorf("failed querying task %d: %w", id, err)
	}
	if project.Valid {
		t.Project = project.String
	}
	return t, nil
}

// GetTasksByStatus retrieves tasks filtered by status.
func (tdb *TaskDB) GetTasksByStatus(status string) ([]task.Task, error) {
	tasks := []task.Task{}
	rows, err := tdb.db.Query("SELECT id, name, project, status, created FROM tasks WHERE status = ? ORDER BY created ASC", status)
	if err != nil {
		return nil, fmt.Errorf("unable to query tasks by status %q: %w", status, err)
	}
	defer rows.Close()
	for rows.Next() {
		var t task.Task
		var project sql.NullString
		err = rows.Scan(&t.ID, &t.Name, &project, &t.Status, &t.Created)
		if err != nil {
			return nil, fmt.Errorf("failed scanning task row (by status): %w", err)
		}
		if project.Valid {
			t.Project = project.String
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating task rows (by status): %w", err)
	}
	return tasks, nil
}
