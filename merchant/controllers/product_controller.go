package controllers

import (
	"GOMISPLUS/models"
	"GOMISPLUS/responses"
	"errors"
	"net/http"
	"strconv"
)

func (server *Server) ProductController(w http.ResponseWriter, r *http.Request) {
	var err error
	prdct := models.Product{}

	if r.Method == "OPTIONS" {
		responses.JSON(w, http.StatusOK, nil)
	}

	if r.Method == "POST" {
		productid := r.FormValue("productid")
		productname := r.FormValue("productname")
		hna := r.FormValue("hna")
		unit := r.FormValue("unit")

		product := models.Product{}

		fHna, _ := strconv.ParseFloat(hna, 64)

		product.ProductId = productid
		product.ProductName = productname
		product.Hna = fHna
		product.Unit = unit

		product.Prepare()

		err = product.Validate()
		if err != nil {
			responses.ERROR(w, http.StatusOK, err)
			return
		}

		created := product.SaveProduct(server.PG)
		if created.Note != nil {
			responses.ERROR(w, http.StatusOK, created.Note)
			return
		}

		responses.JSON(w, http.StatusOK, created)
	}

	if r.Method == "GET" {
		productid := r.FormValue("productid")

		if productid != "" {
			data, err := prdct.FindProductByID(server.PG, productid)
			responses.SUCCESS(w, http.StatusOK, err, data)
		} else {
			data, err := prdct.FindAllProducts(server.PG)
			responses.SUCCESS(w, http.StatusOK, err, data)
		}
	}

	if r.Method == "PUT" {
		productname := r.FormValue("productname")
		hna := r.FormValue("hna")
		unit := r.FormValue("unit")

		productidupdate := r.FormValue("productidupdate")

		fHna, _ := strconv.ParseFloat(hna, 64)

		prdct.ProductName = productname
		prdct.Hna = fHna
		prdct.Unit = unit

		if prdct.ProductName == "" {
			responses.ERROR(w, http.StatusOK, errors.New("required product name"))
			return
		}
		if prdct.Hna == 0 {
			responses.ERROR(w, http.StatusOK, errors.New("required hna"))
			return
		}
		if prdct.Unit == "" {
			responses.ERROR(w, http.StatusOK, errors.New("required unit"))
			return
		}
		if productidupdate == "" {
			responses.ERROR(w, http.StatusOK, errors.New("required product id update"))
			return
		}

		updated := prdct.UpdateAProduct(server.PG, productidupdate)
		if updated.Note != nil {
			responses.ERROR(w, http.StatusOK, updated.Note)
			return
		}

		responses.JSON(w, http.StatusOK, updated)
	}

	if r.Method == "DELETE" {
		productiddelete := r.FormValue("productiddelete")

		if productiddelete == "" {
			responses.ERROR(w, http.StatusOK, errors.New("required product id delete"))
			return
		}

		prdct.ProductId = productiddelete

		deleted := prdct.DeleteAProduct(server.PG)
		if deleted.Note != nil {
			responses.ERROR(w, http.StatusOK, deleted.Note)
			return
		}

		responses.JSON(w, http.StatusOK, deleted)
	}
}
