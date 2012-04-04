package tuxychat

import (
	"appengine"
	"appengine/channel"
	"appengine/memcache"
	"errors"
)

type chatRoom struct {
	Users map[string]bool
}

type command struct {
	Type, Email, Message string
}

func newRoom() *chatRoom {
	return &chatRoom{make(map[string]bool)}
}

func createRoom(c appengine.Context, roomId string) error {
	return memcache.JSON.Set(c, &memcache.Item{
		Key:    roomId,
		Object: newRoom(),
	})
}

func roomExists(c appengine.Context, roomId string) (bool, error) {
	_, err := memcache.JSON.Get(c, roomId, newRoom())
	if err == nil {
		return true, nil
	} else if err == memcache.ErrCacheMiss {
		return false, nil
	}

	return false, err
}

func joinRoom(c appengine.Context, roomId string, email string) (string, error) {
	room := newRoom()
	item, err := memcache.JSON.Get(c, roomId, room)
	if err != nil {
		return "", err
	}

	room.Users[email] = true
	item.Object = room

	if err := memcache.JSON.Set(c, item); err != nil {
		return "", err
	}

	publish(c, roomId, email, "")
	return channel.Create(c, email+roomId)
}

func publish(c appengine.Context, roomId string, email string, message string) error {
	room := newRoom()
	_, err := memcache.JSON.Get(c, roomId, room)
	if err != nil {
		return err
	}

	errs := make([]error, 0, len(room.Users))
	for user, _ := range room.Users {
		if user == email {
			continue
		}

		if message == "" {
			channel.SendJSON(c, user+roomId, command{"join", email, message})
		} else {
			channel.SendJSON(c, user+roomId, command{"msg", email, message})
		}
	}

	if len(errs) > 0 {
		return errors.New("Publishing message failed!")
	}

	return nil
}
