package gist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	UserNameNotSet = fmt.Errorf("username is not set.")
	TokenNotSet    = fmt.Errorf("token is not set.")
)

func getUserInfo() (username, token string, err error) {
	username = os.Getenv("GIT_USERNAME")
	if username == "" {
		err = UserNameNotSet
	}
	token = os.Getenv("GIT_TOKEN")
	if token == "" {
		err = TokenNotSet
	}
	return
}

func request(method, path string, data []byte) (response string, err error) {
	username, token, err := getUserInfo()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(method, "https://api.github.com"+path, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(username, token)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
}

func get(path string) (string, error) {
	return request("GET", path, nil)
}

func post(path string, data []byte) (string, error) {
	return request("POST", path, data)
}

func patch(path string, data []byte) (string, error) {
	return request("PATCH", path, data)
}

func GetList() (result []Overview, err error) {
	if raw, err := get("/gists"); err != nil {
		return nil, err
	} else {
		err = json.Unmarshal([]byte(raw), &result)
		return result, err
	}
}

func GetStarredList() (result []Overview, err error) {
	if raw, err := get("/gists/starred"); err != nil {
		return nil, err
	} else {
		err = json.Unmarshal([]byte(raw), &result)
		return result, err
	}
}

func GetGist(id string) (*Gist, error) {
	var result Gist
	if raw, err := get("/gists/" + id); err != nil {
		return nil, err
	} else {
		err = json.Unmarshal([]byte(raw), &result)
		return &result, err
	}
}

func CreateGist(filename, description, content string) (*Gist, error) {
	newGist := NewGist{
		Description: description,
		Public:      true,
		Files:       map[string]NewGistFile{filename: NewGistFile{Content: content}},
	}
	bytes, err := json.Marshal(newGist)
	if err != nil {
		return nil, err
	}

	if result, err := post("/gists", bytes); err != nil {
		return nil, err
	} else {
		var gist Gist
		err := json.Unmarshal([]byte(result), &gist)
		return &gist, err
	}
}

func UpdateGist(id, filename, description, content string) (*Gist, error) {
	editGist := EditGist{
		Files: map[string]NewGistFile{filename: NewGistFile{Content: content}},
	}
	if description != "" {
		editGist.Description = description
	}
	bytes, err := json.Marshal(editGist)
	if err != nil {
		return nil, err
	}

	if result, err := patch("/gists/"+id, bytes); err != nil {
		return nil, err
	} else {
		var gist Gist
		err := json.Unmarshal([]byte(result), &gist)
		return &gist, err
	}
}
