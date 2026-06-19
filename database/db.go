package database

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type User struct {
	ID       string    `json:"id"`
	Balance  int       `json:"balance"`
	Bank     int       `json:"bank"`
	LastWork time.Time `json:"last_work"`
	Color    string    `json:"color"`
}

var (
	users = make(map[string]User)
	mu    sync.Mutex
	file  = "users.json"
)

func Init() {
	load()
}

func load() {
	data, err := os.ReadFile(file)
	if err != nil {
		return
	}
	json.Unmarshal(data, &users)
}

func save() {
	data, _ := json.MarshalIndent(users, "", "  ")
	os.WriteFile(file, data, 0644)
}

func CreateUser(id string) {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := users[id]; exists {
		return
	}
	users[id] = User{
		ID:       id,
		Balance:  0,
		Bank:     0,
		LastWork: time.Now().Add(-24 * time.Hour),
		Color:    "",
	}
	save()
}

func GetUser(id string) User {
	mu.Lock()
	defer mu.Unlock()
	if u, ok := users[id]; ok {
		return u
	}
	return User{
		ID:       id,
		Balance:  0,
		Bank:     0,
		LastWork: time.Now().Add(-24 * time.Hour),
		Color:    "",
	}
}

func GetAllUsers() []User {
	mu.Lock()
	defer mu.Unlock()
	result := make([]User, 0, len(users))
	for _, u := range users {
		result = append(result, u)
	}
	return result
}

func UpdateBalance(id string, amount int) {
	mu.Lock()
	defer mu.Unlock()
	u := users[id]
	u.Balance += amount
	users[id] = u
	save()
}

func UpdateBank(id string, amount int) {
	mu.Lock()
	defer mu.Unlock()
	u := users[id]
	u.Bank += amount
	users[id] = u
	save()
}

func UpdateWork(id string) {
	mu.Lock()
	defer mu.Unlock()
	u := users[id]
	u.LastWork = time.Now()
	users[id] = u
	save()
}

func SetColor(id string, color string) {
	mu.Lock()
	defer mu.Unlock()
	u := users[id]
	u.Color = color
	users[id] = u
	save()
}
