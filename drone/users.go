package drone

import (
	"fmt"
)

type UserService struct {
	*Client
}

// GET /api/users/{host}/{login}
func (s *UserService) Get(login string) (*User, error) {
	var path = fmt.Sprintf("/api/users/%s", login)
	var user = User{}
	var err = s.run("GET", path, nil, &user)
	return &user, err
}

// GET /api/user
func (s *UserService) GetCurrent() (*User, error) {
	var user = User{}
	var err = s.run("GET", "/api/user", nil, &user)
	return &user, err
}

// POST /api/user/sync
func (s *UserService) Sync() error {
	var err error
	if !s.isServer04 {
		err = s.run("POST", "/api/user/sync", nil, nil)
	}
	return err
}

// POST /api/users/{host}/{login}
func (s *UserService) Create(remote, login string, in interface{}) (*User, error) {
	var path string
	if s.isServer04 {
		path = fmt.Sprintf("/api/users")
	} else {
		path = fmt.Sprintf("/api/users/%s/%s", remote, login)
	}
	var user = User{}
	var err = s.run("POST", path, in, &user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (s *UserService) Patch(login string, in interface{}) (*User, error) {
	var path string

	path = fmt.Sprintf("/api/users/%s", login)

	var user = User{}
	var err = s.run("PATCH", path, in, &user)
	return &user, err
}

// DELETE /api/users/{host}/{login}
func (s *UserService) Delete(remote, login string) error {
	var path string
	if s.isServer04 {
		path = fmt.Sprintf("/api/users/%s", login)
	} else {
		path = fmt.Sprintf("/api/users/%s/%s", remote, login)
	}
	return s.run("DELETE", path, nil, nil)
}

// GET /api/users
func (s *UserService) List() ([]*User, error) {
	var users []*User
	var err = s.run("GET", "/api/users", nil, &users)
	return users, err
}
