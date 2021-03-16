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

//New creates a new DBManager instance and returns it
func New(user, pass, name string) (*DBManager, error) {
	d := DBManager{dbUser: user, dbPass: pass, dbName: name}
	err := d.connectToPSQL()
	return &d, err
}

func (db *DBManager) getPasswordHash(email string) ([]byte, error) {
	result := db.DB.QueryRow(`SELECT password from users WHERE email = $1`, email)

	var hash []byte
	if err := result.Scan(&hash); err != nil {
		return nil, err
	}
	return hash, nil
}

//CheckUserCredentials returns an Error if the supplied credentials do not match any row in the database.
func (db *DBManager) CheckUserCredentials(email, password string) error {
	hash, err := db.getPasswordHash(email)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}

//AddNewUser returns an Error if the user is not successfully added to DB
func (db *DBManager) AddNewUser(email, password string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	_, err = db.DB.Exec(`INSERT INTO users(email,password) VALUES ($1 , $2);`, email, hashedPass)
	if err != nil {
		return err
	}
	return nil
}

//ConnectToPSQL returns an error if the connection to the DB is not successful
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
	_, err := db.DB.Exec(`
		CREATE TABLE IF NOT EXISTS users(
			 uid serial PRIMARY KEY,
			 email VARCHAR (254) UNIQUE NOT NULL,
			 password VARCHAR (60) NOT NULL,
			 created TIMESTAMP NOT NULL default current_timestamp
			 )`,
	)
	return err
}
