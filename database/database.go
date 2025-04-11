package database

import (
	"fmt"
	"github.com/tobischo/gokeepasslib/v3"
	"os"
)

func LoadDatabase(fileName string, password string) (*gokeepasslib.Database, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return nil, fmt.Errorf("Сould not open file %s: %v", fileName, err)
	}
	defer file.Close()

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(password)

	if err := gokeepasslib.NewDecoder(file).Decode(db); err != nil {
		return nil, fmt.Errorf("Failed to decode database: %w", err)
	}

	return db, nil
}
