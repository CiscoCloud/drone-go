package drone

import (
	"errors"
	"fmt"
)

type RepoService struct {
	*Client
}

// GET /api/repos/{host}/{owner}/{name}
func (s *RepoService) Get(host, owner, name string) (*Repo, error) {
	var path string
	if s.isServer04 {
		path = fmt.Sprintf("/api/repos/%s/%s", owner, name)
	} else {
		path = fmt.Sprintf("/api/repos/%s/%s/%s", host, owner, name)
	}
	var repo = Repo{}
	var err = s.run("GET", path, nil, &repo)
	if err == nil {
		return &repo, nil
	} else {
		return nil, err
	}
}

// POST /api/repos/{owner}/{name}
func (s *RepoService) Create(owner, name string) (*Repo, error) {
	if !s.isServer04 {
		return nil, errors.New("No create repos method before Drone 0.4")
	}
	path := fmt.Sprintf("/api/repos/%s/%s?no-activate=true", owner, name)
	var result = Repo{}
	var err = s.run("POST", path, nil, &result)
	if err == nil {
		return &result, nil
	} else {
		return nil, err
	}
}

// PUT /api/repos/{host}/{owner}/{name}
func (s *RepoService) Update(repo *Repo) (*Repo, error) {
	var path, method string
	if s.isServer04 {
		path = fmt.Sprintf("/api/repos/%s/%s", repo.Owner, repo.Name)
		method = "PATCH"
	} else {
		path = fmt.Sprintf("/api/repos/%s/%s/%s", repo.Host, repo.Owner, repo.Name)
		method = "PUT"
	}
	var result = Repo{}
	var err = s.run(method, path, &repo, &result)
	return &result, err
}

// POST /api/repos/{host}/{owner}/{name}
func (s *RepoService) Enable(host, owner, name string) (*Repo, error) {
	var path string
	if s.isServer04 {
		path = fmt.Sprintf("/api/repos/%s/%s", owner, name)
	} else {
		path = fmt.Sprintf("/api/repos/%s/%s/%s", host, owner, name)
	}
	var result = Repo{}
	var err error
	err = s.run("POST", path, nil, &result)
	if err != nil {
		return nil, err
	} else {
		return &result, nil
	}
}

func (s *RepoService) EnableWithActivate(host, owner, name string, activate bool) (*Repo, error) {
	var path string
	if s.isServer04 {
		path = fmt.Sprintf("/api/repos/%s/%s?no-activate=%v", owner, name, !activate)
	} else {
		path = fmt.Sprintf("/api/repos/%s/%s/%s", host, owner, name)
	}
	var result = Repo{}
	var err error
	err = s.run("POST", path, nil, &result)
	if err != nil {
		return nil, err
	} else {
		return &result, err
	}
}

// POST /api/repos/{host}/{owner}/{name}/deactivate
func (s *RepoService) Disable(host, owner, name string) error {
	var path string
	if s.isServer04 {
		return errors.New("No disable function in Drone 0.4")
	} else {
		path = fmt.Sprintf("/api/repos/%s/%s/%s/deactivate", host, owner, name)
	}
	return s.run("POST", path, nil, nil)
}

// DELETE /api/repos/{host}/{owner}/{name}?remove=true
func (s *RepoService) Delete(host, owner, name string) error {
	var path string
	if s.isServer04 {
		path = fmt.Sprintf("/api/repos/%s/%s", owner, name)
	} else {
		path = fmt.Sprintf("/api/repos/%s/%s/%s", host, owner, name)
	}
	return s.run("DELETE", path, nil, nil)
}

// PUT /api/repos/{host}/{owner}/{name}
func (s *RepoService) SetKey(host, owner, name, pub, priv string) error {
	var path string
	if s.isServer04 {
		path = fmt.Sprintf("/api/repos/%s/%s", owner, name)
	} else {
		path = fmt.Sprintf("/api/repos/%s/%s/%s", host, owner, name)
	}
	var in = struct {
		PublicKey  string `json:"public_key"`
		PrivateKey string `json:"private_key"`
	}{pub, priv}
	return s.run("PUT", path, &in, nil)
}

// PUT /api/repos/{host}/{owner}/{name}
func (s *RepoService) SetParams(host, owner, name, params interface{}) error {
	var path, method string
	if s.isServer04 {
		path = fmt.Sprintf("/api/repos/%s/%s", owner, name)
		method = "PATCH"
	} else {
		path = fmt.Sprintf("/api/repos/%s/%s/%s", host, owner, name)
		method = "PUT"
	}
	var in = struct {
		Params interface{} `json:"params"`
	}{params}
	return s.run(method, path, &in, nil)
}

// GET /api/user/repos
func (s *RepoService) List() ([]*Repo, error) {
	var repos []*Repo
	var err = s.run("GET", "/api/user/repos", nil, &repos)
	return repos, err
}
