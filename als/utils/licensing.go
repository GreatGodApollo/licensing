package utils

import (
	"database/sql"
	"errors"
	"github.com/GreatGodApollo/als/crypto"
	"github.com/GreatGodApollo/als/database"
	"github.com/spf13/viper"
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {

	s := make([]byte, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func generateLicenseString() string {
	return RandomString(4) + "-" + RandomString(4) + "-" + RandomString(4)
}

func GenerateEncryptedLicense(db *sql.DB, product, email string) ([]byte, error) {
	key := generateLicenseString()

	exist, err := database.CheckLicenseExist(db, key)
	if err != nil {
		return nil, err
	}
	if !exist {
		query, err := db.Prepare("insert licenses SET license_key=?, product=?, email=?")
		if err != nil {
			return nil, err
		}
		_, err = query.Exec(key, product, email)
		if err != nil {
			return nil, err
		}
		defer query.Close()
		return crypto.Encrypt([]byte(viper.GetString("crypt.key")), []byte(key))
	}
	return nil, errors.New("license already exists")
}

func DecryptLicense(encrypted []byte) ([]byte, error) {
	return crypto.Decrypt([]byte(viper.GetString("crypt.key")), encrypted)
}
