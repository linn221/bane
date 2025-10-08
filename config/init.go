package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var _BASE_DIR string

// get the base directory, load environment variables
func init() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	_BASE_DIR = dir
	environmentPath := filepath.Join(dir, ".env")
	err = godotenv.Load(environmentPath)
	if err != nil {
		log.Fatal(err)
	}
}
func GetBaseDir() string {
	return _BASE_DIR
}
