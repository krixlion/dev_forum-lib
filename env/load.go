package env

import (
	"fmt"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

// Load assumes the ".env" file is located in the root directory.
// It returns an error if it cannot find a file named ".env" in the project root dir.
func Load(projectDir string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to load .env: %w", err)
		}
	}()

	re, err := regexp.Compile(`^(.*` + projectDir + `)`)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	rootPath := re.Find([]byte(cwd))

	return godotenv.Load(string(rootPath) + `/.env`)
}
