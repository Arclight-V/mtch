package repository

type DBKind string

const (
	Postgres DBKind = "postgres"
	Memory   DBKind = "memory"
)
