package model

import "time"

type UserSession struct {
    ID        string    `gorm:"primaryKey" json:"id"`
    UserID    string    `json:"user_id"`
    JWTToken  string    `gorm:"unique" json:"jwt_token"`
    CreatedAt time.Time `json:"created_at"`
    ExpiresAt time.Time `json:"expires_at"`
}
