package variables

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	Bottoken string = os.Getenv("BOTTOKEN")
	DB       string = os.Getenv("DB")
)
