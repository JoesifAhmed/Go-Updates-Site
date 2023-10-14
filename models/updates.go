package models

import (
	"fmt"
	"strconv"
	"time"
)

type Update struct {
	id int64
}

func NewUpdate(userId int64, body string) (*Update, error) {
	id, err := Client.Incr("update:next-id").Result()
	updateTime := time.Now().Local().Format("2006-Jan-02 3:4 pm")

	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("update:%d", id)
	pipe := Client.Pipeline()
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "user_id", userId)
	pipe.HSet(key, "body", body)
	pipe.HSet(key, "update_time", updateTime)
	pipe.LPush("updates", id)
	pipe.LPush(fmt.Sprintf("user:%d:updates", userId), id)
	_, err = pipe.Exec()
	if err != nil {
		return nil, err
	}
	return &Update{id}, nil
}
func (update *Update) GetBody() (string, error) {
	key := fmt.Sprintf("update:%d", update.id)
	return Client.HGet(key, "body").Result()
}
func (update *Update) GetUser() (*User, error) {
	key := fmt.Sprintf("update:%d", update.id)
	userId, err := Client.HGet(key, "user_id").Int64()
	if err != nil {
		return nil, err
	}
	return GetUserById(userId)
}
func (update *Update) GetTime() (string, error) {
	key := fmt.Sprintf("update:%d", update.id)
	return Client.HGet(key, "update_time").Result()

}

func queryUpdates(key string) ([]*Update, error) {
	updateIds, err := Client.LRange(key, 0, 10).Result()
	if err != nil {
		return nil, err
	}
	updates := make([]*Update, len(updateIds))
	for i, strid := range updateIds {
		id, err := strconv.Atoi(strid)
		if err != nil {
			return nil, err
		}
		updates[i] = &Update{int64(id)}
	}
	return updates, nil
}
func GetAllUpdates() ([]*Update, error) {
	return queryUpdates("updates")
}
func GetUpdates(userId int64) ([]*Update, error) {
	key := fmt.Sprintf("user:%d:updates", userId)
	return queryUpdates(key)
}

func PostUpdate(userId int64, body string) error {
	_, err := NewUpdate(userId, body)
	return err
}
