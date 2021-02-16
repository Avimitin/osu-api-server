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
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Date
}

type Data struct {
	PlayCount string `json:"play_count"`
	Rank      string `json:"rank"`
	PP        string `json:"pp"`
	Acc       string `json:"acc"`
	PlayTime  string `json:"play_time"`
}

type Date struct {
	Recent    Data `json:"recent"`
	Yesterday Data `json:"yesterday"`
}
