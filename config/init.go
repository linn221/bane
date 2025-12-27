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
		// Don't fatal on missing .env file (useful for tests)
		// Try loading from current working directory as fallback
		if cwd, err := os.Getwd(); err == nil {
			envPath := filepath.Join(cwd, ".env")
			if err := godotenv.Load(envPath); err != nil {
				// .env file is optional, just log a warning
				log.Printf("Warning: Could not load .env file from %s or %s: %v", environmentPath, envPath, err)
			}
		}
	}
}
func GetBaseDir() string {
	return _BASE_DIR
}
