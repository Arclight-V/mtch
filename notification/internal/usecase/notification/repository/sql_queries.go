package repository

const insertIssueSql = `INSERT INTO verify_codes_issued (user_id, code, expires_at, purpose, attempts, max_attempts)
						VALUES ($1, $2, $3, $4, $5, $6)
						ON CONFLICT (user_id, purpose) DO UPDATE
						SET code 		= EXCLUDED.code,
						expires_at  	= EXCLUDED.expires_at,
						attempts    	= 0,
						max_attempts 	= EXCLUDED.max_attempts;`
