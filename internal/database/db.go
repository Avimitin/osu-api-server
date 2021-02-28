package database

import (
	"errors"
	"fmt"
	"time"
)

//OsuDB contain sql.DB field
type OsuDB struct {
	UsersData PlayerDataStore
}

// Connect return database connection by given DSN
func Connect(driver string, dsn string) (*OsuDB, error) {
	switch driver {
	case "mysql":
		store, err := NewMySQLStore(dsn)
		if err != nil {
			return nil, fmt.Errorf("connect %s:%v", dsn, err)
		}
		// Set limit
		if store != nil {
			store.db.SetConnMaxLifetime(time.Minute * 3)
			store.db.SetMaxOpenConns(10)
			store.db.SetMaxIdleConns(10)
		}
		err = store.db.Ping()
		if err != nil {
			return nil, fmt.Errorf("connect db: %v", err)
		}
		return &OsuDB{store}, nil
	case "redis":
		store := NewRedisDataStore()
		return &OsuDB{store}, nil
	}
	return nil, errors.New("unsupport database driver")
}

func (db *OsuDB) CheckUserDataStoreHealth() error {
	if db.UsersData == nil {
		return errors.New("userdata has not yet initialized")
	}
	return db.UsersData.CheckHealth()
}

// GetUserRecent return user data with given name
func (db *OsuDB) GetUserRecent(username string) (*User, error) {
	return db.UsersData.GetPlayer(username)
}

// GetUserYtd return a user's yesterday data with given name
func (db *OsuDB) GetUserYtd(username string) (*User, error) {
	return db.UsersData.GetPlayerOld(username)
}

// InsertNewUser insert user data into database
func (db *OsuDB) InsertNewUser(
	userID string, username string, pc string, rank string, pp string, acc string, playtime string,
) error {
	return db.UsersData.AddPlayer(
		User{
			UserID:   userID,
			Username: username,
			Date: Date{
				Recent: Data{
					PlayCount: pc,
					Rank:      rank,
					PP:        pp,
					Acc:       acc,
					PlayTime:  playtime,
				},
			},
		},
	)
}

// UpdateUser update user recent data with given data
func (db *OsuDB) UpdateUser(
	username string, pc string, rank string, pp string, acc string, playtime string,
) error {
	return db.UsersData.Update(
		User{
			Username: username,
			Date: Date{
				Recent: Data{
					PlayCount: pc,
					Rank:      rank,
					PP:        pp,
					Acc:       acc,
					PlayTime:  playtime,
				},
			},
		})
}

func (db *OsuDB) UpdateUserYtd(
	username string, pc string, rank string, pp string, acc string, playtime string,
) error {
	return db.UsersData.UpdateOld(
		User{
			Username: username,
			Date: Date{
				Yesterday: Data{
					PlayCount: pc,
					Rank:      rank,
					PP:        pp,
					Acc:       acc,
					PlayTime:  playtime,
				},
			},
		})
}
