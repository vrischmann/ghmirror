package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/vrischmann/ghmirror/internal/config"
)

func makeDB(conf *config.Postgres) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", conf.Host, conf.Port, conf.User, conf.Password, conf.Dbname, conf.SSLMode)
	return sql.Open("postgres", dsn)
}
