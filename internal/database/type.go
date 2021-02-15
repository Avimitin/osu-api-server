package database

type PlayerDataStore interface {
	AddPlayer(User) error
	CheckHealth() error
	GetPlayer(string) (*User, error)
	GetPlayerOld(string) (*User, error)
	Update(User) error
	UpdateOld(User) error
}
