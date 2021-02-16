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
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	PlayCount    string `json:"play_count"`
	Rank         string `json:"rank"`
	PP           string `json:"pp"`
	Acc          string `json:"acc"`
	TotalPlay    string `json:"total_play"`
	PcYtd        string `json:"play_count_ytd"`
	RankYtd      string `json:"rank_ytd"`
	PpYtd        string `json:"pp_ytd"`
	AccYtd       string `json:"acc_ytd"`
	TotalPlayYtd string `json:"ttp_ytd"`
}
