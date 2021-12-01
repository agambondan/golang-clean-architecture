module go-blog-api

// +heroku goVersion go1.16

go 1.16

require (
	github.com/badoux/checkmail v1.2.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/contrib v0.0.0-20201101042839-6a891bf89f19
	github.com/gin-gonic/gin v1.7.2
	github.com/go-redis/redis/v7 v7.4.0
	github.com/google/uuid v1.2.0
	github.com/joho/godotenv v1.3.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/oauth2 v0.0.0-20210805134026-6f1e6394065a
	google.golang.org/api v0.30.0
	gorm.io/driver/postgres v1.1.2
	gorm.io/gorm v1.21.15
)
