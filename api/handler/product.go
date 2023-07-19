package handler

import (
	"app/models"
	"app/pkg/helper"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (h *handler) Product(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		h.CreateProduct(w, r)
	case "GET":
		var (
			values = r.URL.Query()
			method = values.Get("method")
		)

		if method == "GET_LIST" {
			h.GetListProduct(w, r)
		} else if method == "GET" {
			h.GetByIdProduct(w, r)
		}
	case "PUT":
		h.UpdateProduct(w, r)
	case "DELETE":
		h.DeleteProduct(w, r)
	}

}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {

	var createProduct models.CreateProduct

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.handlerResponse(w, "error while read body: "+err.Error(), http.StatusBadRequest, nil)
		return
	}

	err = json.Unmarshal(body, &createProduct)
	if err != nil {
		h.handlerResponse(w, "error while unmarshal body: "+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	id, err := h.strg.Product().CreateProduct(&createProduct)
	if err != nil {
		h.handlerResponse(w, "error while storage product create:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	resp, err := h.strg.Product().GetByIDProduct(&models.ProductPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(w, "error while storage product get by id:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(w, "Success", http.StatusOK, resp)
}

func (h *handler) GetByIdProduct(w http.ResponseWriter, r *http.Request) {

	var id string = r.URL.Query().Get("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(w, "invalid id uuid", http.StatusBadRequest, nil)
		return
	}

	resp, err := h.strg.Product().GetByIDProduct(&models.ProductPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(w, "error while storage product get by id:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(w, "Success", http.StatusOK, resp)
}

func (h *handler) GetListProduct(w http.ResponseWriter, r *http.Request) {

	var (
		offsetStr = r.URL.Query().Get("offset")
		limitStr  = r.URL.Query().Get("limit")
		search    = r.URL.Query().Get("search")
	)

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		h.handlerResponse(w, "error while offset: "+err.Error(), http.StatusBadRequest, nil)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.handlerResponse(w, "error while limit: "+err.Error(), http.StatusBadRequest, nil)
		return
	}

	resp, err := h.strg.Product().GetListProduct(&models.ProductGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: search,
	})
	if err != nil {
		h.handlerResponse(w, "error while storage product get list:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(w, "Success", http.StatusOK, resp)
}

func (h *handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	var (
		resp = models.UpdateProduct{}
	)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.handlerResponse(w, "error while read body: "+err.Error(), http.StatusBadRequest, nil)
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		h.handlerResponse(w, "error while unmarshal body: "+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	prodModel, err := h.strg.Product().UpdateProduct(&resp)
	if err != nil {
		h.handlerResponse(w, "error while storage product update:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(w, "Succes", http.StatusOK, prodModel)
}

func (h *handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var (
		resp = models.ProductPrimaryKey{}
	)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.handlerResponse(w, "error while read body: "+err.Error(), http.StatusBadRequest, nil)
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		h.handlerResponse(w, "error while unmarshal body: "+err.Error(), http.StatusInternalServerError, nil)
		return
	}
	err = h.strg.Product().DeleteProduct(&resp)
	if err != nil {
		h.handlerResponse(w, "error while storage product delete:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(err.Error()))
}
