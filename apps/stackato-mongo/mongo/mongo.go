package mongo

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/cloudfoundry-community/go-cfenv"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/hpcloud/catalog-templates/apps/stackato-mongo/common"
)

const collection = "testusers"

// User db Model
type User struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Username string        `json:"user"`
	Email    string        `json:"email"`
	Password string        `json:"password"`
}

var database string

func dbConnection() (*mgo.Session, error) {
	appEnv, err := cfenv.Current()
	if err != nil {
		errstr := "Failed to get the current cf application environment: " + err.Error()
		return nil, errors.New(errstr)
	}
	svcList, err := appEnv.Services.WithTag("stackato-mongo")
	if err != nil {
		errstr := "Failed to get the application environment service: " + err.Error()
		return nil, errors.New(errstr)
	}
	svc := svcList[0]
	uri, ok := svc.CredentialString("uri")
	if !ok {
		errstr := "failed to get uri connection string from env"
		return nil, errors.New(errstr)
	}
	database, ok = svc.CredentialString("db")
	if !ok {
		errstr := "failed to get database from env"
		return nil, errors.New(errstr)
	}

	common.Logger.Debug("RUNNING IN CF MODE WITH MONGO")
	session, err := mgo.Dial(uri)
	if err != nil {
		errstr := "failed to connect database: " + err.Error()
		return nil, errors.New(errstr)
	}
	return session, nil
}

func closeConnection(session *mgo.Session) {
	session.Close()
}

func init() {
	session, err := dbConnection()
	if err != nil {
		common.Logger.Fatal(err)
	}
	defer closeConnection(session)

	// make a clean collection in database
	err = session.DB(database).C(collection).DropCollection()
	if err != nil {
		common.Logger.Error("Collection/namespace already clean: ", err)
	}
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
	session, err := dbConnection()
	if err != nil {
		common.Logger.Fatal(err)
		return nil, err
	}
	defer closeConnection(session)

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
	c := session.DB(database).C(collection)
	err = c.Insert(&user)
	if err != nil {
		common.Logger.Error(err)
		return nil, err
	}
	return &user, nil
}

// Read gets the data from the database and returns it
func Read() ([]User, error) {
	session, err := dbConnection()
	if err != nil {
		common.Logger.Fatal(err)
		return nil, err
	}
	defer closeConnection(session)

	c := session.DB(database).C(collection)
	iter := c.Find(nil).Iter()
	var users []User
	err = iter.All(&users)
	if err != nil {
		common.Logger.Fatal(err)
		return nil, err
	}
	if users == nil {
		return []User{}, nil
	}
	return users, nil
}
