package services

import (
	"github.com/tradesim/model"
	"math/rand"
	"time"
	"fmt"
)

// session key => session object (username and expiration time)
var SessionStore map[string]model.Session
var duration int

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func InitializeSessions() {
	SessionStore = make(map[string]model.Session)
	duration = int(20 * time.Minute) / int(time.Millisecond)
	go KillExpiredSessions()
}

// creates a new session and adds it to the session map, or replaces it if the user already has a session
func CreateSession(username string) string {
	// first delete any existing sessions for this user
	for key, val := range SessionStore {
		if val.Username == username {
			delete(SessionStore, key)
		}
	}

	// next, create new session
	sessionKey := RandomString(32)
	var session model.Session
	session.Username = username
	session.Expiration = int(time.Now().UnixNano() / int64(time.Millisecond)) + duration
	SessionStore[sessionKey] = session
	return sessionKey
}

func GetSession(sessionKey string) (model.Session, error) {
	if val, ok := SessionStore[sessionKey]; ok {
		username := SessionStore[sessionKey].Username
		SessionStore[sessionKey] = model.Session{
			Username: username,
			Expiration: int(time.Now().UnixNano() / int64(time.Millisecond)) + duration,
		}
	    return val, nil
	} else {
		return model.Session{}, fmt.Errorf("Session not found")
	}
}

func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func DeleteSessionByKey(key string) {
	delete(SessionStore, key)
}

// runs in background periodically destroying expired sessions
func KillExpiredSessions() {
	for {
		time.Sleep(time.Second)
		now := int(time.Now().UnixNano() / int64(time.Millisecond))
		for key, val := range SessionStore {
			if val.Expiration < now {
				fmt.Println("Deleting session with key " + key)
				fmt.Println("Username: " + val.Username)
				delete(SessionStore, key)
			}
		}
	}
}











