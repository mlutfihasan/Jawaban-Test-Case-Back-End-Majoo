package controllers

import (
	"GOMISPLUS/models"
	"GOMISPLUS/responses"
	"errors"
	"net/http"
)

type CrudResult struct {
	Status string `json:"status"`
	Note   error  `json:"note"`
}

func (server *Server) LoginController(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method == "OPTIONS" {
		responses.JSON(w, http.StatusOK, nil)
	}

	if r.Method == "GET" {
		nik := r.FormValue("nik")
		password := r.FormValue("password")

		user := models.User{}

		user.Nik = nik
		user.Password = password

		err = user.AuthenticateLogin(server.PG)
		if err != nil {
			responses.ERROR(w, http.StatusOK, errors.New("wrong nik or password"))
			return
		}

		result := CrudResult{
			Status: "1",
			Note:   nil,
		}

		responses.JSON(w, http.StatusCreated, result)
	}
}
