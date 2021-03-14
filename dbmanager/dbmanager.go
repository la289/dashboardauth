package dbmanager

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq" //db driver for postgres
	"golang.org/x/crypto/bcrypt"

)

var ErrUserNonexistant = errors.New("User does not exist")

type DBManager struct {
	dbUser, dbPass, dbName string
	DB                     *sql.DB
}

func New(user, pass, name string) (DBManager, error) {
	d := DBManager{dbUser: user, dbPass: pass, dbName: name}
	err := d.connectToPSQL()
	return d, err
}

func (db *DBManager) getPasswordHash(email string) ([]byte, error) {
	rows, err := db.DB.Query(`SELECT password from users WHERE email = $1`, email)
	if err != nil {
		return nil, err
	}

	hasNextRow := rows.Next()
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if !hasNextRow {
		return nil, ErrUserNonexistant
	}
	var hash []byte
	if err = rows.Scan(&hash); err != nil {
		return nil, err
	}
	if len(hash) != 60 {
		return nil, ErrUserNonexistant
	}

	return hash, nil
}


// returns an Error if the credentials are not valid
func (db *DBManager) CheckUserCredentials(email, password string) error {
	hash, err := db.getPasswordHash(email)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}


// Returns err if user is not added to DB
func (db *DBManager) AddNewUser(email, password string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	_, err = db.DB.Query(`INSERT INTO users(email,password) VALUES ($1 , $2);`, email, hashedPass)
	if err != nil {
		return err
	}
	return nil
}

// Returns an error if the connection is not made
func (db *DBManager) connectToPSQL() error {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		db.dbUser, db.dbPass, db.dbName)
	psql, err := sql.Open("postgres", dbinfo)

	if err != nil {
		return err
	}
	db.DB = psql
	// only creates the table if it doesn't already exist
	// Typically this database would exist on its own outside of docker
	err = db.initSchemaUsers()
	if err != nil {
		return err
	}

	return nil
}

func (db *DBManager) initSchemaUsers() error {
	_, err := db.DB.Query(`
		CREATE TABLE IF NOT EXISTS users(
			 uid serial PRIMARY KEY,
			 email VARCHAR (254) UNIQUE NOT NULL,
			 password VARCHAR (60) NOT NULL,
			 created TIMESTAMP NOT NULL default current_timestamp
			 )`,
	)
	return err
}
