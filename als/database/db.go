package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/GreatGodApollo/als/crypto"
	"github.com/GreatGodApollo/als/models"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
)

func Setup() (*sql.DB, error) {
	return sql.Open("mysql", getConnectionString())
}

func getConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.name"))
}

func CheckLicenseExist(db *sql.DB, key string) (bool, error) {

	var scanned string
	err := db.QueryRow("select license_key from licenses where license_key = ?", key).Scan(&scanned)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func CheckLicenseValid(db *sql.DB, key string) (bool, bool, error) {

	exist, err := CheckLicenseExist(db, key)
	if err != nil {
		return false, false, err
	}
	if !exist {
		return false, false, nil
	}

	var valid bool
	err = db.QueryRow("select valid from licenses where license_key=?", key).Scan(&valid)
	if err != nil {
		return false, false, err
	} else {
		return exist, valid, nil
	}
}

func CheckLicenseValidProduct(db *sql.DB, key, product string) (bool, bool, error) {

	exist, err := CheckLicenseExist(db, key)
	if err != nil {
		return false, false, err
	}
	if !exist {
		return false, false, nil
	}

	var valid bool
	var prodScanned string
	err = db.QueryRow("select valid, product from licenses where license_key=?", key).Scan(&valid, &prodScanned)
	if !valid {
		return exist, valid, nil
	}

	if err != nil {
		return false, false, err
	} else {
		if product == prodScanned && valid {
			return exist, valid, nil
		} else if valid {
			return exist, valid, errors.New("incorrect product")
		} else {
			return exist, false, nil
		}
	}
}

func InvalidateLicense(db *sql.DB, key string) (bool, error) {
	exist, valid, err := CheckLicenseValid(db, key)
	if err != nil {
		return false, err
	}

	if exist && valid {
		query, err := db.Prepare("update licenses set valid=? where license_key=?")
		if err != nil {
			return false, err
		}
		_, err = query.Exec(false, key)
		defer query.Close()
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	} else if exist {
		return false, errors.New("license already invalid")
	} else {
		return false, errors.New("license nonexistent")
	}
}

func GetWholeRecord(db *sql.DB, key string) (models.License, error) {
	exist, err := CheckLicenseExist(db, key)
	if err != nil {
		return models.License{}, err
	}
	if exist {
		var licObj models.License

		err = db.QueryRow("select * from licenses where license_key = ?", key).Scan(&licObj.Id,
			&licObj.LicenseKey,
			&licObj.Product,
			&licObj.Email,
			&licObj.Valid)
		if err != nil {
			return models.License{}, err
		}
		return licObj, nil
	} else {
		return models.License{}, errors.New("license nonexistent")
	}
}

func GetAllValidRecords(db *sql.DB, product string) (models.Licenses, error) {
	rows, err := db.Query("select * from licenses where valid = 1 and product = ?", product)
	if err != nil {
		return models.Licenses{}, err
	}

	got := models.Licenses{}.Licenses
	for rows.Next() {
		var r models.License
		err = rows.Scan(&r.Id,
			&r.LicenseKey,
			&r.Product,
			&r.Email,
			&r.Valid)
		if err != nil {
			return models.Licenses{}, err
		}

		encr, err := crypto.Encrypt([]byte(viper.GetString("crypt.key")), []byte(r.LicenseKey))
		if err != nil {
			return models.Licenses{}, err
		}
		r.LicenseKey = crypto.EncodeBase64(encr)

		got = append(got, r)
	}

	licensObj := models.Licenses{}
	licensObj.Licenses = got

	return licensObj, nil

}
