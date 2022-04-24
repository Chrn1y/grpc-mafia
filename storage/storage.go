package storage

import (
	"github.com/rs/xid"
)

type Gender int

const (
	Male   Gender = 1
	Female Gender = 2
)

type User struct {
	ID     string
	Name   string
	Avatar []byte
	Gender Gender
	Email  string
}

type Storage struct {
	users map[string]*User
}

func (s *Storage) Insert(inp *User) string { // post
	id := xid.New().String()
	s.users[id] = inp
	return id
}

func (s *Storage) Delete(id string) { // delete
	if _, ok := s.users[id]; ok {
		delete(s.users, id)
	}
}

func (s *Storage) Get(id string) (*User, bool) { // get
	user, ok := s.users[id]
	return user, ok
}

func (s *Storage) Validate(name string) bool {
	for _, u := range s.users {
		if name == u.Name {
			return true
		}
	}
	return false
}

func (s *Storage) Update(id string, user *User) bool {
	if _, ok := s.users[id]; ok {
		s.users[id] = user
		return true
	}
	return false
}

func (s *Storage) List() []*User { // get
	res := make([]*User, 0, len(s.users))
	for _, u := range s.users {
		res = append(res, u)
	}
	return res
}

func New() *Storage {
	return &Storage{
		users: make(map[string]*User),
	}
}
