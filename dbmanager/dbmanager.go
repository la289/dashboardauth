package dbmanager

import(
	"database/sql"
	_ "github.com/lib/pq" //db driver for postgres
	"fmt"
	"iotdashboard/utils"
	"errors"
)

//////// QUERIES
const (
	userExistsQuery = `SELECT uid from users WHERE email = $1`
	getUserInfoQuery      = `SELECT * from users WHERE email = $1`
	getUserPassQuery      = `SELECT password from users WHERE email = $1`
	createUserQuery = `INSERT INTO users(uid,email,password) VALUES (DEFAULT, $1 , $2);`
)

//////// ERRORS
var ErrUserNonexistant = errors.New("User does not exist")
var ErrUserExists 	   = errors.New("User already exists")


type DBManager struct{
	dbUser, dbPass, dbName string
	DB *sql.DB
}

func New(user, pass, name string) (DBManager, error) {
	d := DBManager{dbUser: user, dbPass: pass, dbName: name}
	err := d.connectToPSQL()
	return d, err
}

// VerifyUserExists returns an error if the user does not exist
func (db *DBManager) VerifyUserExists(email string) error {
	rows, err := db.DB.Query(userExistsQuery, email)
	if err != nil {
		return err
	}

	if !rows.Next(){
		return  ErrUserNonexistant
	}
	return nil
}


func (db *DBManager) GetUserPassHash(email string) ([]byte, error) {
	err := db.VerifyUserExists(email)
	if err != nil {
		return nil, err
	}

	rows, err := db.DB.Query(getUserPassQuery, email)
	if err != nil{
		return nil, err
	}

	rows.Next()
	var hash []byte
	rows.Scan(&hash)
	return hash, nil
}


// Returns err if user is not added to DB
func (db *DBManager) AddNewUser(email, password string) error {
	err := db.VerifyUserExists(email)
	if err == nil { //TODO this needs to check for the exact error, otherwise return the error. Need errors to be better defined
		return ErrUserExists
	} else if err != ErrUserNonexistant {
		return err
	}

	hashedPass, err := utils.GeneratePassHash(password)
	if err != nil{
		return err
	}

	_, err = db.DB.Query(createUserQuery, email, hashedPass)
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
	err = db.createAccountsTable()
	if err != nil {
		return err
	}

	return nil
}

func (db *DBManager) createAccountsTable() error {
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


