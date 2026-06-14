package users

import (
	"errors"
	"fmt"
	"log/slog"
	"net/mail"
	"time"
)

var ErrNoResultsFound = errors.New("No results found")

type User struct {
	FirstName string
	LastName  string
	Email     mail.Address
}

type Manager struct {
	users []User
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) AddUser(firstname, lastname, email string) error {
	parsedEmail, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("Invalid email address: %s", email)
	}

	if firstname == "" {
		return fmt.Errorf("Invalid first name: %q", firstname)
	} else if lastname == "" {
		return fmt.Errorf("Invalid last name: %q", lastname)
	}

	existingUser, err := m.GetUserByName(firstname, lastname)
	if err != nil && !errors.Is(err, ErrNoResultsFound) {
		return fmt.Errorf("Error checking if user exists: %v", err)
	}
	if existingUser != nil {
		return fmt.Errorf("User already exists")
	}

	newUser := User{firstname, lastname, *parsedEmail}
	m.users = append(m.users, newUser)
	return nil
}

func (m *Manager) GetUserByName(firstName, lastName string) (*User, error) {
	for i, user := range m.users {
		if user.FirstName == firstName && user.LastName == lastName {
			result := m.users[i]
			return &result, nil
		}
	}

	return nil, ErrNoResultsFound
}

func (m *Manager) Shutdown() {
	slog.Info("Manager shutting down")
	time.Sleep(2 * time.Second)
	slog.Info("Manager shutdown complete")
}
