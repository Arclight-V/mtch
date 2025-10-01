package repository

const insertIssueSql = `INSERT INTO verify_tokens_issued (jti, user_id, expires_at) VALUES ($1, $2, $3)`
