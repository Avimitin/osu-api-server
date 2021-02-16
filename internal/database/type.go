package database

type PlayerDataStore interface {
	AddPlayer(User) error
	CheckHealth() error
	GetPlayer(string) (*User, error)
	GetPlayerOld(string) (*User, error)
	Update(User) error
	UpdateOld(User) error
}

// User type contain user field
type User struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	PlayCount Date   `json:"play_count"`
	Rank      Date   `json:"rank"`
	PP        Date   `json:"pp"`
	Acc       Date   `json:"acc"`
	PlayTime  Date   `json:"play_time"`
}

type Date struct {
	Recent    string `json:"recent"`
	Yesterday string `json:"yesterday"`
}
