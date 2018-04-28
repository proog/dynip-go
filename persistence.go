package main

import (
	"encoding/json"
	"log"

	"github.com/gomodule/redigo/redis"
)

// An IPInfo contains information about a client
type IPInfo struct {
	IP      string
	Updated string
}

// A Persistence contains data necessary to save and load ip infos
type Persistence struct {
	redis  redis.Conn
	memory map[string]IPInfo
}

// NewPersistence makes a new persistence
func NewPersistence(redisAddress string) Persistence {
	conn, error := redis.Dial("tcp", redisAddress)
	memory := make(map[string]IPInfo)

	if error != nil {
		log.Println("Redis connection unavailable, falling back to in-memory store")
		conn = nil
	}

	return Persistence{
		redis:  conn,
		memory: memory,
	}
}

// Load loads an ipinfo from persistence
func Load(p Persistence, name string) (IPInfo, bool) {
	if p.redis == nil {
		value, ok := p.memory[name]
		return value, ok
	}

	value, _ := redis.Bytes(p.redis.Do("GET", name))

	if value == nil {
		return IPInfo{}, false
	}

	unmarshalled := IPInfo{}
	json.Unmarshal(value, &unmarshalled)

	return unmarshalled, true
}

// Save saves an ipinfo to persistence
func Save(p Persistence, name string, info IPInfo) {
	if p.redis == nil {
		p.memory[name] = info
		return
	}

	value, _ := json.Marshal(info)
	p.redis.Do("SET", name, value)
}
