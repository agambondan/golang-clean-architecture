package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"golang-youtube-api/config"
	"log"
)

type Repositories struct {
	Role         RoleRepository
	User         UserRepository
	Category     CategoryRepository
	Post         PostRepository
	PostCategory PostCategoryRepository
	db           *sql.DB
}

func NewRepositories(configure config.Configuration) (*Repositories, error) {
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
		Role:         NewRoleRepository(db),
		User:         NewUserRepository(db),
		Category:     NewCategoryRepository(db),
		Post:         NewPostRepository(db),
		PostCategory: NewPostCategoryRepository(db),
		db:           db,
	}, nil
}

// Close closes the  database connection
func (s *Repositories) Close() error {
	return s.db.Close()
}

// Seeder This migrate all tables
func (s *Repositories) Seeder() error {
	var err error
	var result sql.Result
	result, err = s.db.Exec("CREATE TABLE IF NOT EXISTS users (uuid uuid PRIMARY KEY, first_name VARCHAR(55) not null, last_name VARCHAR(55) not null, email VARCHAR(55) unique not null, " +
		"phone_number VARCHAR(15) not null, username VARCHAR(55) unique not null, password VARCHAR(255) not null, photo_profile VARCHAR(255), role_id int not null, " +
		"instagram VARCHAR(255), facebook VARCHAR(255), twitter VARCHAR(255), linkedin VARCHAR(255), created_at timestamp, updated_at timestamp, deleted_at timestamp)")
	// log.Println(result, err)
	result, err = s.db.Exec("CREATE TABLE IF NOT EXISTS roles (id serial PRIMARY KEY, name VARCHAR(15) not null unique, " +
		"created_at timestamp, updated_at timestamp, deleted_at timestamp)")
	// log.Println(result, err)
	result, err = s.db.Exec("CREATE TABLE IF NOT EXISTS posts (id serial PRIMARY KEY, title VARCHAR(100) not null unique, description text not null, thumbnail VARCHAR(255), user_uuid uuid not null, " +
		"created_at timestamp, updated_at timestamp, deleted_at timestamp)")
	// log.Println(result, err)
	result, err = s.db.Exec("CREATE TABLE IF NOT EXISTS categories (id serial PRIMARY KEY, name VARCHAR(35) not null unique, " +
		"created_at timestamp, updated_at timestamp, deleted_at timestamp)")
	// log.Println(result, err)
	result, err = s.db.Exec("CREATE TABLE IF NOT EXISTS post_categories " +
		"(post_id serial not null, category_id int not null)")
	// log.Println(result, err)
	if err != nil {
		log.Println(result, err)
	}
	return err
}

func (s *Repositories) AddForeignKey() error {
	var err error
	return err
}
