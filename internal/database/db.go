package database

import (
	"database/sql"
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
	var store PlayerDataStore
	switch driver {
	case "mysql":
		store = &MySQLDataStore{}
		db, err := sql.Open(driver, dsn)
		if err != nil {
			return nil, fmt.Errorf("connect %s:%v", dsn, err)
		}
		// Set limit
		if db != nil {
			db.SetConnMaxLifetime(time.Minute * 3)
			db.SetMaxOpenConns(10)
			db.SetMaxIdleConns(10)
		}
	default:
		return nil, errors.New("unsupport database driver")
	}

	return &OsuDB{store}, nil
}

func (db *OsuDB) CheckUserDataStoreHealth() error {
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
	userID string, username string, pc string, rank string, pp string, acc string, total_play string,
) error {
	return db.UsersData.AddPlayer(
		User{
			UserID:    userID,
			Username:  username,
			PlayCount: Date{Recent: pc},
			Rank:      Date{Recent: rank},
			PP:        Date{Recent: pp},
			Acc:       Date{Recent: acc},
			PlayTime:  Date{Recent: total_play},
		},
	)
}

// UpdateUser update user data with given data
func (db *OsuDB) UpdateUser(
	username string, pc string, rank string, pp string, acc string, total_play string,
) error {
	return db.UsersData.Update(
		User{
			Username:  username,
			PlayCount: Date{Recent: pc},
			Rank:      Date{Recent: rank},
			PP:        Date{Recent: pp},
			Acc:       Date{Recent: acc},
			PlayTime:  Date{Recent: total_play},
		})
}

func (db *OsuDB) UpdateUserYtd(
	username string, pc string, rank string, pp string, acc string, total_play string,
) error {
	return db.UsersData.UpdateOld(
		User{
			Username:  username,
			PlayCount: Date{Yesterday: pc},
			Rank:      Date{Yesterday: rank},
			PP:        Date{Yesterday: pp},
			Acc:       Date{Yesterday: acc},
			PlayTime:  Date{Yesterday: total_play},
		})
}
