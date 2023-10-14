package models

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidLogin      = errors.New("incorrect password")
	ErrUserAlreadyExists = errors.New("User Already Exists")
)

type User struct {
	id int64
}

func NewUser(username string, hash []byte) (*User, error) {
	exists, err := Client.HExists("user:by-username", username).Result()
	if exists {
		return nil, ErrUserAlreadyExists
	}
	if err != nil {
		return nil, err
	}
	id, err := Client.Incr("user:next-id").Result()

	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("user:%d", id)
	pipe := Client.Pipeline()
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "username", username)
	pipe.HSet(key, "hash", hash)
	pipe.HSet("user:by-username", username, id)
	_, eror := pipe.Exec()
	if eror != nil {
		return nil, eror
	}
	return &User{id}, eror
}

func (user *User) GetUsername() (string, error) {
	key := fmt.Sprintf("user:%d", user.id)
	return Client.HGet(key, "username").Result()
}
func (user *User) GetId() (int64, error) {
	return user.id, nil
}
func (user *User) GetHash() ([]byte, error) {
	key := fmt.Sprintf("user:%d", user.id)
	return Client.HGet(key, "hash").Bytes()
}

func GetUserById(id int64) (*User, error) {
	return &User{id}, nil
}

func GetUserByUsername(username string) (*User, error) {
	id, err := Client.HGet("user:by-username", username).Int64()
	if err == redis.Nil {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	return GetUserById(id)
}

func (user *User) Authnitcate(password string) error {
	hash, err := user.GetHash()
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrInvalidLogin
	}
	return err
}

func AuthnticateUser(username, password string) (*User, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, user.Authnitcate(password)
}

func RegiserUser(username, password string) error {
	cost := bcrypt.DefaultCost

	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	if err != nil {
		return err
	}
	_, err = NewUser(username, hash)
	return err
}
