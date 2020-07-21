package models

import (
	"api/pkg/setting"
	"github.com/go-redis/redis"
	"log"
	"time"
)

var (
	Rdb *redis.Client
)

func init() {
	// 获取redis 配置
	r, err := setting.Cfg.GetSection("redis")
	if err != nil {
		log.Fatal(2, "Fail to get section 'redis': %v", err)
	}

	ins, _ := r.Key("DB").Int()
	ConnectRedis(r.Key("ADDRESS").String(),
				r.Key("PASSWD").String(),
				ins)

	// 每次服务启动清除本机ip地址缓存
	Rdb.Del(ServerLocalRunIpKey)
}

func ConnectRedis(addr string, passwd string, db int){
	Rdb = redis.NewClient(&redis.Options{
		Addr 		: addr,
		Password	: passwd,
		DB       	: db,
	})
}

func DelKey(key string) {
	Rdb.Del(key).Val()
}

func GetValByKey(key string) interface{} {
	return  Rdb.Get(key).Val()
}

func SetValByKey(key string, val interface{}, expiration time.Duration) error{
	_, err :=Rdb.Set(key, val, expiration).Result()

	return  err
}

func SetValBySetKey(key string, val interface{}) error{
	_, err := Rdb.SAdd(key, val).Result()

	return  err
}

func CheckMemberByKey(key string, val interface{}) bool{
	isMember, _ := Rdb.SIsMember(key, val).Result()
	return isMember
}




