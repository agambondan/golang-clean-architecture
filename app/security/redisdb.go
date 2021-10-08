package security

import "github.com/go-redis/redis/v7"


type RedisService struct {
	Auth   Interface
	Client *redis.Client
}

func NewRedisDB(host, port, password string) (*RedisService, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})
	return &RedisService{
		Auth:   NewAuthRedis(redisClient),
		Client: redisClient,
	}, nil
}