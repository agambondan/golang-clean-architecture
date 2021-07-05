package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"golang-youtube-api/config"
	"log"
)

type Repositories struct {
	User UserRepository
	db   *sql.DB
}

func NewRepositories(configure config.Config) (*Repositories, error) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		configure.DBHost, configure.DBPort, configure.DBUser, configure.DBName, configure.DBPassword)
	db, _ := sql.Open(configure.DBDriver, DBURL)
	err := db.Ping()
	if err != nil {
		fmt.Printf("Cannot connect to %s database ", configure.DBDriver)
		log.Fatalln("\nThis is the error:", err)
		return nil, err
	} else {
		fmt.Printf("We are connected to the %s database with url %s\n", configure.DBDriver, DBURL)
	}
	return &Repositories{
		User: NewUserRepository(db),
		db:   db,
	}, nil
}

// Close closes the  database connection
func (s *Repositories) Close() error {
	return s.db.Close()
}

// Seeder This migrate all tables
func (s *Repositories) Seeder() error {
	var err error
	return err
}

func (s *Repositories) AddForeignKey() error {
	var err error
	return err
}
