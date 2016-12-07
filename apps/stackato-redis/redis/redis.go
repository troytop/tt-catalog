package redis

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/cloudfoundry-community/go-cfenv"
	"gopkg.in/redis.v4"

	"github.com/hpcloud/catalog-templates/apps/stackato-redis/common"
)

func dbConnection() (*redis.Client, error) {
	appEnv, err := cfenv.Current()
	if err != nil {
		errstr := "Failed to get the current cf application environment: " + err.Error()
		return nil, errors.New(errstr)
	}
	svcList, err := appEnv.Services.WithTag("stackato-redis")
	if err != nil {
		errstr := "Failed to get the application environment service: " + err.Error()
		return nil, errors.New(errstr)
	}
	svc := svcList[0]
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
	port, ok := svc.CredentialString("port")
	if !ok {
		errstr := "failed to get port from env"
		return nil, errors.New(errstr)
	}

	common.Logger.Debug("Building connection to redis")
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})
	pong, err := client.Ping().Result()
	if err != nil || pong != "PONG" {
		errstr := "failed to connect database: " + err.Error()
		return nil, errors.New(errstr)
	}
	common.Logger.Debug("Connected to redis")
	return client, nil
}

func closeConnection(db *redis.Client) {
	db.Close()
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

// Entry is a data entry in redis database
type Entry struct {
	Key   string
	Value interface{}
}

// Read gets the data from the database and returns it
func Read() ([]Entry, error) {
	client, err := dbConnection()
	if err != nil {
		return nil, err
	}
	defer closeConnection(client)
	keys := client.Keys("*").Val()
	rd := []Entry{}
	for _, key := range keys {
		n := client.Get(key)
		rd = append(rd, Entry{key, n.Val()})
	}
	return rd, nil
}

// Write generates some random data and inserts it into the database
func Write() (*Entry, error) {
	client, err := dbConnection()
	if err != nil {
		return nil, err
	}
	defer closeConnection(client)
	key, err := generateString(5, usercharacters)
	if err != nil {
		return nil, err
	}
	value, err := generateString(10, passwordcharacters)
	if err != nil {
		return nil, err
	}
	n := client.Set(key, value, 0)
	n.Err()
	if err != nil {
		return nil, err
	}
	common.Logger.Debug("wrote data to db: " + n.Val())
	return &Entry{key, value}, nil
}
