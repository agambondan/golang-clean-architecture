package repository

import (
	"fmt"
	"go-blog-api/app/config"
	"go-blog-api/app/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
)

type Repositories struct {
	Role     RoleRepository
	User     UserRepository
	Category CategoryRepository
	Article  ArticleRepository
	db       *gorm.DB
}

func NewRepositories(configure config.Configuration) (*Repositories, error) {
	logLevel := logger.Info

	switch os.Getenv("ENVIRONMENT") {
	case "staging":
		logLevel = logger.Error
	case "production":
		logLevel = logger.Silent
	}
	gormConfig := gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   os.Getenv("DB_TABLE_PREFIX"),
			SingularTable: true,
		},
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta", configure.DBHost, configure.DBPort, configure.DBUser, configure.DBPassword, configure.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gormConfig)
	if err != nil {
		fmt.Printf("Cannot connect to %s database url %s\n", configure.DBDriver, dsn)
		dsn = fmt.Sprint("postgres://cgxtqgoobyvbwk:47465f1dd068148279716e2788dc252a6bb85339d5b2d635d2b4557b5c7e2627@ec2-34-204-128-77.compute-1.amazonaws.com:5432/d594pchn88flmk")
		db, err = gorm.Open(postgres.Open(dsn), &gormConfig)
		if err != nil {
			log.Fatalln(err)
		}
		if db.Error != nil {
			log.Fatalln(err)
		}
	}
	if err != nil {
		fmt.Printf("Cannot connect to %s database url %s", configure.DBDriver, dsn)
		log.Println("\nThis is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database with url %s\n", configure.DBDriver, dsn)
	}
	return &Repositories{
		Role:     NewRoleRepository(db),
		User:     NewUserRepository(db),
		Category: NewCategoryRepository(db),
		Article:  NewArticleRepository(db),
		db:       db,
	}, nil
}

// Close closes the  database connection
func (s *Repositories) Close() error {
	db, _ := s.db.DB()
	return db.Close()
}

func (s *Repositories) Migrations() error {
	err := s.db.AutoMigrate(migrations.ModelMigrations...)
	err = s.db.Migrator().DropTable("schema_migration")
	return err
}

// Seeder This migrate all tables
func (s *Repositories) Seeder() error {
	return nil
}

func (s *Repositories) AddForeignKey() error {
	var err error
	return err
}
