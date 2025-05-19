package variables

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	Bottoken string = os.Getenv("BOTTOKEN")
)
