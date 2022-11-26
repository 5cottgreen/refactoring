package internal

import (
	"encoding/json"
	"errors"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const store = `users.json`

var UserNotFound = errors.New("user not found")

type errResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

// Renders error response payload
func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// Initializes invalid request renderer
func errInvalidRequest(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

// Checks if service is working
func CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(time.Now().String()))
}

// Fetches users
func FetchUsers(w http.ResponseWriter, r *http.Request) {
	f, err := ioutil.ReadFile(store)
	if err != nil {
		log.Fatalf("'%s' faild to read file", err.Error())
	}

	s := UserStore{}
	if err := json.Unmarshal(f, &s); err != nil {
		log.Fatalf("'%s' faild to unmarshal", err.Error())
	}

	render.JSON(w, r, s.List)
}

// Creates users
func CreateUser(w http.ResponseWriter, r *http.Request) {
	request := CreateUserRequest{}

	if err := render.Bind(r, &request); err != nil {
		if err := render.Render(w, r, errInvalidRequest(err)); err != nil {
			log.Fatalf("'%s' faild to render", err.Error())
		}
		return
	}

	f, err := ioutil.ReadFile(store)
	if err != nil {
		log.Fatalf("'%s' faild to read file", err.Error())
	}
	s := UserStore{}
	if err := json.Unmarshal(f, &s); err != nil {
		log.Fatalf("'%s' faild to unmarshal", err.Error())
	}

	s.Increment++
	u := User{
		CreatedAt: time.Now(),
		Name:      request.Name,
		Email:     request.Email,
	}

	id := strconv.Itoa(s.Increment)
	s.List[id] = u

	b, err := json.Marshal(&s)
	if err != nil {
		log.Fatalf("'%s' faild to marshal", err.Error())
	}

	if err := ioutil.WriteFile(store, b, fs.ModePerm); err != nil {
		log.Fatalf("'%s' faild to write into file", err.Error())
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"user_id": id,
	})
}

// Receives user by param
func GetUser(w http.ResponseWriter, r *http.Request) {
	f, err := ioutil.ReadFile(store)
	if err != nil {
		log.Fatalf("'%s' faild to read file", err.Error())
	}
	s := UserStore{}
	if err := json.Unmarshal(f, &s); err != nil {
		log.Fatalf("'%s' faild to unmarshal", err.Error())
	}

	id := chi.URLParam(r, "id")

	render.JSON(w, r, s.List[id])
}

// Updates user by param
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	request := UpdateUserRequest{}

	if err := render.Bind(r, &request); err != nil {
		if err := render.Render(w, r, errInvalidRequest(err)); err != nil {
			log.Fatalf("'%s' faild to render", err.Error())
		}
		return
	}

	f, err := ioutil.ReadFile(store)
	if err != nil {
		log.Fatalf("'%s' faild to read file", err.Error())
	}
	s := UserStore{}
	if err := json.Unmarshal(f, &s); err != nil {
		log.Fatalf("'%s' faild to unmarshal", err.Error())
	}

	id := chi.URLParam(r, "id")

	if _, ok := s.List[id]; !ok {
		if err = render.Render(w, r, errInvalidRequest(UserNotFound)); err != nil {
			log.Fatalf("'%s' faild to render", err.Error())
		}
		return
	}

	u := s.List[id]
	u.Name = request.Name
	s.List[id] = u

	b, err := json.Marshal(&s)
	if err != nil {
		log.Fatalf("'%s' faild to marshal", err.Error())
	}

	if err := ioutil.WriteFile(store, b, fs.ModePerm); err != nil {
		log.Fatalf("'%s' faild to write into file", err.Error())
	}

	render.Status(r, http.StatusNoContent)
}

// Deletes user by param
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	f, err := ioutil.ReadFile(store)
	if err != nil {
		log.Fatalf("'%s' faild to read file", err.Error())
	}
	s := UserStore{}
	if err := json.Unmarshal(f, &s); err != nil {
		log.Fatalf("'%s' faild to unmarshal", err.Error())
	}

	id := chi.URLParam(r, "id")

	if _, ok := s.List[id]; !ok {
		if err := render.Render(w, r, errInvalidRequest(UserNotFound)); err != nil {
			log.Fatalf("'%s' faild to render", err.Error())
		}
		return
	}

	delete(s.List, id)

	b, err := json.Marshal(&s)
	if err != nil {
		log.Fatalf("'%s' faild to marshal", err.Error())
	}

	if err := ioutil.WriteFile(store, b, fs.ModePerm); err != nil {
		log.Fatalf("'%s' faild to write into file", err.Error())
	}

	render.Status(r, http.StatusNoContent)
}
