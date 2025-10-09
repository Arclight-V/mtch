package repository

const (
	createPendingUserQuery = `INSERT INTO users (email, password)
								VALUES ($1, $2)
								RETURNING user_id, email, role, created_at, updated_at`

	findByEmailQuery = `SELECT user_id, email, first_name, last_name, role, avatar, password, created_at, updated_at FROM users WHERE email = $1`
)
