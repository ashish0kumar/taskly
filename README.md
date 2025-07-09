<h1 align="center">Taskly</h1>

<p align="center">
Command-line task manager
</p>

| **List View**                  | **Kanban View**                   |
| ------------------------------ | --------------------------------- |
| ![List View](assets/table.png) | ![Kanban View](assets/kanban.png) |

## Features

- **Task Management:** Add, delete, update, and list tasks.
- **Project Organization:** Assign tasks to specific projects for better
  organization.
- **SQLite Database Integration:** Efficient storage and retrieval of tasks.
- **Styled Output:** Stylish table and Kanban layouts using Lip Gloss.
- **Kanban Board Interface:** Visualize tasks as a Kanban board with Bubble Tea.

## Installation

**Prerequisites:**

- [Go](https://golang.org/doc/install) (version 1.16 or later recommended)

1. Clone the repository:

   ```bash
   git clone https://github.com/ashish0kumar/taskly.git
   cd taskly
   ```

2. Install the application using `go install`:

   ```bash
   go install .
   ```

This command will compile the `main.go` file and place the resulting executable
(`taskly`) in your `GOPATH/bin` directory. Ensure that your `GOPATH/bin` is
included in your system's PATH environment variable so you can run `taskly` from
anywhere in your terminal.

Alternatively, you can build the executable manually:

```bash
go build -o taskly .
# Then move the 'taskly' executable to a directory in your PATH
# For example: mv taskly /usr/local/bin/
```

## Usage

Use Taskly commands from your terminal to manage tasks. The main command is
`taskly`, followed by subcommands to perform specific actions.

### Commands

- **Add a Task:** Add a new task, optionally specifying a project name

  ```bash
  taskly add "Task Name" -p "Project Name"
  ```

- **Delete a Task:** Delete a task by its unique ID:

  ```bash
  taskly delete <ID>
  ```

- **Update a Task:** Update a task's name, project, or status:

  ```bash
  taskly update <ID> -n "New Task Name" -p "New Project Name" -s <status>
  ```

  _Status options:_
  - `0` for "todo"
  - `1` for "in progress"
  - `2` for "done"

- **List All Tasks:** List all stored tasks in a table format:

  ```bash
  taskly list
  ```

- **View Kanban Board:** Display tasks in a Kanban board layout. Tasks are
  categorized into `todo`, `in progress` and `done` columns:

  ```bash
  taskly kanban
  ```

- **View Database Path:** Locate the database file where tasks are stored:

  ```bash
  taskly where
  ```

## Examples

1. **Adding a Task with Project Name**

```bash
taskly add "Design Homepage" -p "Website Redesign"
```

2. **Updating a Task's Status**

```bash
taskly update 1 -s 1
```

3. **Listing All Tasks**

```bash
taskly list
```

This command shows tasks in a formatted table with columns for ID, Name,
Project, Status, and Creation Date.

4. **Viewing the Kanban Board**

```bash
taskly kanban
```

This displays tasks in a Kanban layout, categorized by status.

### Data Storage

Taskly uses a SQLite database to persist tasks. The database is stored in an
XDG-compliant directory (typically `$HOME/.local/share/tasks.db`). This
structure enables easy backup and integration across systems.

## Dependencies

- [Cobra](https://github.com/spf13/cobra): CLI command framework.
- [Bubble Tea](https://github.com/charmbracelet/bubbletea): TUI framework for
  the Kanban board.
- [Lip Gloss](https://github.com/charmbracelet/lipgloss): Used for stylish table
  layouts.
- [SQLite](https://github.com/mattn/go-sqlite3): Lightweight, serverless SQL
  database.

## Contributions

Contributions are welcome! Please fork the repository, create a new branch, and
submit a pull request.

<br>

<p align="center">
	<img src="https://raw.githubusercontent.com/catppuccin/catppuccin/main/assets/footers/gray0_ctp_on_line.svg?sanitize=true" />
</p>

<p align="center">
        <i><code>&copy 2025-present <a href="https://github.com/ashish0kumar">Ashish Kumar</a></code></i>
</p>

<div align="center">
<a href="https://github.com/ashish0kumar/taskly/blob/main/LICENSE"><img src="https://img.shields.io/github/license/ashish0kumar/taskly?style=for-the-badge&color=CBA6F7&logoColor=cdd6f4&labelColor=302D41" alt="LICENSE"></a>&nbsp;&nbsp;
</div>
