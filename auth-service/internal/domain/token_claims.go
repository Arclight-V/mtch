package domain

import "time"

type TokenClaims struct {
	UserId string
	Role   string
	Exp    time.Time
}
