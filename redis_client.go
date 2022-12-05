package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
)

var INTERNAL_ERROR = "Internal error. Please try again later."

var ctx = context.Background()
var rAddr = "localhost:6379"
var rPass = "localhost:6379"

var rdb = redis.NewClient(&redis.Options{
	Addr:     rAddr,
	Password: rPass, // no password set
	DB:       0,     // use default DB
})

var rdbNames = redis.NewClient(&redis.Options{ //name: user map
	Addr:     rAddr,
	Password: rPass, // no password set
	DB:       1,     // use default DB
})

func SetUser(id string, name string) (string, bool) {
	if err := rdb.Set(ctx, id, name, 0).Err(); err != nil {
		return INTERNAL_ERROR, false
	}
	if err := rdbNames.Set(ctx, name, id, 0).Err(); err != nil {
		rdb.Del(ctx, id, name)
		return INTERNAL_ERROR, false
	}

	return "", true
}

func GetUser(id string) (string, bool) {
	name, err := rdb.Get(ctx, id).Result()
	if err != nil {
		fmt.Printf("%e, %s, %s", err, name, id)
		return "", false
	}

	return name, true
}

func UsernameExists(name string) bool {
	_, err := rdbNames.Get(ctx, name).Result()
	if err != nil {
		fmt.Printf("%e", err)
		return false
	}

	return true
}

// TODO: Cleanup
//func ExampleClient() {
//	err := rdb.Set(ctx, "key", "value", 0).Err()
//	if err != nil {
//		panic(err)
//	}
//
//	val, err := rdb.Get(ctx, "key").Result()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("key", val)
//
//	val2, err := rdb.Get(ctx, "key2").Result()
//	if err == redis.Nil {
//		fmt.Println("key2 does not exist")
//	} else if err != nil {
//		panic(err)
//	} else {
//		fmt.Println("key2", val2)
//	}
//	// Output: key value
//	// key2 does not exist
//}
