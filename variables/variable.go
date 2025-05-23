package variables

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	Bottoken    string = os.Getenv("BOTTOKEN")
	Mongodb_uri string = os.Getenv("MONGODB_URI")
	Db_name     string = os.Getenv("DB_NAME")
)
