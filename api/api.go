package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dkpeakbil/taskserver/domain"
	"github.com/dkpeakbil/taskserver/usecase"
	"net/http"
)

type Api struct {
	addr  string
	ucase usecase.UseCase
}

func NewApi(addr string, ucase usecase.UseCase) (*Api, error) {
	return &Api{
		addr:  addr,
		ucase: ucase,
	}, nil
}

func (a *Api) Run() error {
	http.HandleFunc("/register", a.handleRegisterRequest)
	http.HandleFunc("/auth", a.handleAuthRequest)
	return http.ListenAndServe(a.addr, nil)
}

func (a *Api) handleRegisterRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.Form.Get("username") == "" {
		a.responseBadRequest(w, errors.New("missing username parameter"))
		return
	}

	if r.Form.Get("password") == "" {
		a.responseBadRequest(w, errors.New("missing password parameter"))
		return
	}

	request := &domain.RegisterRequest{
		Username: r.Form.Get("username"),
		Password: r.Form.Get("password"),
	}

	response := a.ucase.Register(request)
	a.responseOK(w, response)
}

func (a *Api) handleAuthRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.Form.Get("username") == "" {
		a.responseBadRequest(w, errors.New("missing username parameter"))
		return
	}

	if r.Form.Get("password") == "" {
		a.responseBadRequest(w, errors.New("missing password parameter"))
		return
	}

	request := &domain.AuthRequest{
		Username: r.Form.Get("username"),
		Password: r.Form.Get("password"),
	}

	response := a.ucase.Auth(request)
	a.responseOK(w, response)
}

func (a *Api) responseBadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = fmt.Fprint(w, err.Error())
}

func (a *Api) responseOK(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(res)
	_, _ = fmt.Fprint(w, string(j))
}
