package repository

const (
	createPendingUserQuery = `INSERT INTO users (user_id, first_name, last_name, contact, phone, email, password, date_birthday, gender, role)
								VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
								RETURNING user_id, first_name, last_name, contact, phone, email, password, deta_birthday, gender, role, created_at, updated_at`

	findByEmailQuery = `SELECT user_id, email, first_name, last_name, role, avatar, password, created_at, updated_at FROM users WHERE email = $1`
)
