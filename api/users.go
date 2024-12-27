package users

import (
	"fmt"

	"github.com/google/uuid"
)

type ID string

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

type UserWithID struct {
	ID ID
	User
}

type UsersApplication struct {
	data map[ID]User
}

func NewApplication() *UsersApplication {
	return &UsersApplication{
		data: make(map[ID]User),
	}
}

func (ua *UsersApplication) FindAll() []UserWithID {
	var users []UserWithID = make([]UserWithID, 0, len(ua.data))
	for id, user := range ua.data {
		users = append(users, UserWithID{
			ID:   id,
			User: user,
		})
	}
	return users
}

func (ua *UsersApplication) FindById(id ID) UserWithID {
	user, ok := ua.data[id]
	if !ok {
		return UserWithID{}
	}
	return UserWithID{
		ID:   id,
		User: user,
	}
}

func (ua *UsersApplication) Insert(u User) UserWithID {
	newUUID := uuid.New().String()

	id := ID(newUUID)

	ua.data[id] = u

	userWithID := UserWithID{
		ID:   id,
		User: u,
	}

	return userWithID
}

func (ua *UsersApplication) Update(id ID, u User) (UserWithID, error) {
	if _, ok := ua.data[id]; !ok {
		return UserWithID{}, fmt.Errorf("user with id %s not found", id)
	}

	ua.data[id] = u
	return UserWithID{
		ID:   id,
		User: u,
	}, nil
}

func (ua *UsersApplication) Delete(id ID) (UserWithID, error) {
	user, ok := ua.data[id]
	if !ok {
		return UserWithID{}, fmt.Errorf("user with id %s not found", id)
	}

	delete(ua.data, id)
	return UserWithID{
		ID:   id,
		User: user,
	}, nil
}
