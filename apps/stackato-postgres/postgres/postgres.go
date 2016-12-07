package postgres

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // needed for gorm

	"github.com/hpcloud/catalog-templates/apps/stackato-postgres/common"
)

// User db Model
type User struct {
	gorm.Model
	Username string `gorm:"size:255" json:"user"`
	Email    string `gorm:"type:varchar(100);unique_index" json:"email"`
	Password string `json:"password"`
}

func dbConnection() (*gorm.DB, error) {
	appEnv, err := cfenv.Current()
	if err != nil {
		errstr := "Failed to get the current cf application environment: " + err.Error()
		return nil, errors.New(errstr)
	}
	svcList, err := appEnv.Services.WithTag("stackato-postgres")
	if err != nil {
		errstr := "Failed to get the application environment service: " + err.Error()
		return nil, errors.New(errstr)
	}
	svc := svcList[0]
	user, ok := svc.CredentialString("user")
	if !ok {
		errstr := "failed to get user from env"
		return nil, errors.New(errstr)
	}
	password, ok := svc.CredentialString("password")
	if !ok {
		errstr := "failed to get password from env"
		return nil, errors.New(errstr)
	}
	host, ok := svc.CredentialString("host")
	if !ok {
		errstr := "failed to get host from env"
		return nil, errors.New(errstr)
	}
	database, ok := svc.CredentialString("database")
	if !ok {
		errstr := "failed to get database from env"
		return nil, errors.New(errstr)
	}
	common.Logger.Debug("RUNNING IN CF MODE WITH POSTGRES")
	connectionString := "host=" + host + " user=" + user + " dbname=" + database + " sslmode=disable password=" + password
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		errstr := "failed to connect database: " + err.Error()
		return nil, errors.New(errstr)
	}
	db.LogMode(true)
	return db, nil
}

func closeConnection(db *gorm.DB) {
	db.Close()
}

func init() {
	db, err := dbConnection()
	if err != nil {
		common.Logger.Fatal(err)
	}
	defer closeConnection(db)

	// Migrate the schema
	db.AutoMigrate(&User{})
}

func generateString(length int, characters string) (string, error) {
	b := make([]byte, length)
	max := big.NewInt(int64(len(characters)))
	for i := range b {
		var c byte
		rint, err := rand.Int(rand.Reader, max)
		if err != nil {
			common.Logger.Error(err)
			return "", errors.New("Unable to generate a string. Error : " + err.Error())
		}
		c = characters[rint.Int64()]
		b[i] = c
	}
	common.Logger.Debug("generated string: ", string(b))
	return string(b), nil
}

const usercharacters = "abcdefghijklmnopqrstuvwxyz"
const passwordcharacters = `abcdefghijklmnopqrstuvwxyz1234567890`

// Write generates a random user
func Write() (*User, error) {
	db, err := dbConnection()
	if err != nil {
		common.Logger.Fatal(err)
		return nil, err
	}
	defer closeConnection(db)
	username, err := generateString(10, usercharacters)
	if err != nil {
		common.Logger.Error("Error generating a username", err)
		return nil, err
	}
	password, err := generateString(10, passwordcharacters)
	if err != nil {
		common.Logger.Error("Error generating a password", err)
		return nil, err
	}
	var user = User{
		Username: username,
		Password: password,
		Email:    username + "@gmail.com",
	}
	common.Logger.Debug("Creating User: ", user)
	common.Logger.Debug("has user been created?: ", db.NewRecord(user))
	db.Create(&user)
	common.Logger.Debug("has user been created?: ", db.NewRecord(user))
	return &user, nil
}

// Read gets the data from the database and returns it
func Read() ([]User, error) {
	db, err := dbConnection()
	if err != nil {
		common.Logger.Fatal(err)
		return nil, err
	}
	defer closeConnection(db)

	var users []User
	db.Find(&users)
	common.Logger.Debug("found users: ", users)

	return users, nil
}
