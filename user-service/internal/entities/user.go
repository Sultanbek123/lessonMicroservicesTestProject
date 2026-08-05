package entities

type TokenClaims struct {
	UserID string
	Role   string
}

type User struct {
	ID           string
	Username     string
	PasswordHash string
	Role         string
}
