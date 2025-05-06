package task

import (
	"fmt"
	"time"
)

// Status represents the state of a task.
type Status int

// Defines the possible task statuses.
const (
	Todo Status = iota
	InProgress
	Done
)

// String returns the string representation of a Status.
func (s Status) String() string {
	statuses := [...]string{"todo", "in progress", "done"}
	if s < 0 || int(s) >= len(statuses) {
		return "unknown"
	}
	return statuses[s]
}

// Task represents a single task item. Exported for use in other packages.
type Task struct {
	ID      uint
	Name    string
	Project string // Use string, handle NULL in DB layer scan
	Status  string // Store as string representation from Status enum
	Created time.Time
}

// list.Item implementation for Bubble Tea lists

func (t Task) FilterValue() string { return t.Name }
func (t Task) Title() string       { return t.Name }
func (t Task) Description() string {
	if t.Project != "" {
		return fmt.Sprintf("Project: %s", t.Project)
	}
	return ""
}

// kancli.Status implementation (Exported methods on exported Status type)
// These methods allow the Status enum to be used with kancli columns.

func (s Status) Next() int {
	if s == Done {
		return int(Todo)
	}
	return int(s + 1)
}

func (s Status) Prev() int {
	if s == Todo {
		return int(Done)
	}
	return int(s - 1)
}

func (s Status) Int() int {
	return int(s)
}

// Helper function to get all defined Status values. Exported
func AllStatuses() []Status {
	return []Status{Todo, InProgress, Done}
}
