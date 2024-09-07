package config

import (
	"github.com/RichardHoa/blog-aggerator/internal/database"
)

type ApiConfig struct {
	DB *database.Queries
}