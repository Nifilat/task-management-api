package utils

// Resources and actions used in permission checks (they match the rows seeded in the permissions table).
const (
	ResourceTask    = "task"
	ResourceUser    = "user"
	ResourceProfile = "profile"

	ActionRead   = "read"
	ActionWrite  = "write"
	ActionDelete = "delete"
)
