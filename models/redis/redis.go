package redis

import (
	"log"

	"github.com/DLTcollab/vehicle-data-explorer/models/config"
	"github.com/go-redis/redis"
)

func ConnectRedisDB() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.ConfigHandler.GetString("REDIS_HOST"),
		Password: config.ConfigHandler.GetString("REDIS_PASSWORD"),
		DB:       config.ConfigHandler.GetInt("REDIS_DB"),
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		log.Print("Can't ping redis server")
	}

	return rdb
}

func New() *RedisStore {
	redisClient := ConnectRedisDB()

	return &RedisStore{
		DB: redisClient,
	}
}

func (db *RedisStore) Set(key string, value string) error {
	err := db.DB.Set(key, value, 0).Err()
	if err != nil {
		log.Panic(err)
	}
	return err
}

func (db *RedisStore) Get(key string) (string, error) {
	val, err := db.DB.Get(key).Result()
	if err != nil {
		log.Panic(err)
	}
	return val, err
}

func (db *RedisStore) Has(key string) (bool, error) {
	_, err := db.DB.Get(key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		log.Panic(err)
		return false, err
	}
	return true, nil
}

func (db *RedisStore) Delete(key string) error {
	_, err := db.DB.Do("del", key).Result()
	if err != nil {
		log.Panic(err)
	}

	return err
}
