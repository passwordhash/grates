package domain

import "time"

type Session struct {
	RefreshToken string `json:"refreshToken" redis:"refreshToken"`
	//ExpiresAt    time.Time `json:"expiresAt" redis:"expiresAt"`
	TTL time.Duration `json:"TTL" redis:"TTL"`
}
