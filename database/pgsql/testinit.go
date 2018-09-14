package pgsql

// Initialize tests
const (
	insert`INSERT INTO namespace (name, version_format) VALUES ($1, $2)`
)
