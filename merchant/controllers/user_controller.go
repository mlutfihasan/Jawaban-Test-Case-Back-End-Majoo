package controllers

import (
	"GOMISPLUS/models"
	"GOMISPLUS/responses"
	"errors"
	"net/http"
)

func (server *Server) UserController(w http.ResponseWriter, r *http.Request) {
	var err error
	usr := models.User{}

	if r.Method == "OPTIONS" {
		responses.JSON(w, http.StatusOK, nil)
	}

	if r.Method == "POST" {
		nik := r.FormValue("nik")
		username := r.FormValue("username")
		phone := r.FormValue("phone")
		email := r.FormValue("email")
		password := r.FormValue("password")

		user := models.User{}

		user.Nik = nik
		user.Username = username
		user.Phone = phone
		user.Email = email
		user.Password = password

		user.Prepare()

		err = user.Validate()
		if err != nil {
			responses.ERROR(w, http.StatusOK, err)
			return
		}

		err = user.BeforeSave()
		if err != nil {
			responses.ERROR(w, http.StatusOK, err)
			return
		}

		created := user.SaveUser(server.PG)
		if created.Note != nil {
			responses.ERROR(w, http.StatusOK, created.Note)
			return
		}

		responses.JSON(w, http.StatusOK, created)
	}

	if r.Method == "GET" {
		nik := r.FormValue("nik")

		if nik != "" {
			data, err := usr.FindUserByID(server.PG, nik)
			responses.SUCCESS(w, http.StatusOK, err, data)
		} else {
			data, err := usr.FindAllUsers(server.PG)
			responses.SUCCESS(w, http.StatusOK, err, data)
		}
	}

	if r.Method == "PUT" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		phone := r.FormValue("phone")

		nikupdate := r.FormValue("nikupdate")

		usr.Username = username
		usr.Email = email
		usr.Phone = phone

		if usr.Username == "" {
			responses.ERROR(w, http.StatusOK, errors.New("required username"))
			return
		}
		if usr.Email == "" {
			responses.ERROR(w, http.StatusOK, errors.New("required email"))
			return
		}
		if usr.Phone == "" {
			responses.ERROR(w, http.StatusOK, errors.New("required phone"))
			return
		}
		if nikupdate == "" {
			responses.ERROR(w, http.StatusOK, errors.New("required nik update"))
			return
		}

		updated := usr.UpdateAUser(server.PG, nikupdate)
		if updated.Note != nil {
			responses.ERROR(w, http.StatusOK, updated.Note)
			return
		}

		responses.JSON(w, http.StatusOK, updated)
	}

	if r.Method == "DELETE" {
		nikdelete := r.FormValue("nikdelete")

		if nikdelete == "" {
			responses.ERROR(w, http.StatusOK, errors.New("required nik delete"))
			return
		}

		usr.Nik = nikdelete

		deleted := usr.DeleteAUser(server.PG)
		if deleted.Note != nil {
			responses.ERROR(w, http.StatusOK, deleted.Note)
			return
		}

		responses.JSON(w, http.StatusOK, deleted)
	}
}
