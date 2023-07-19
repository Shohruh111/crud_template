package handler

import (
	"app/models"
	"app/pkg/helper"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (h *handler) Category(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		h.CreateCategory(w, r)
	case "GET":
		var (
			values = r.URL.Query()
			method = values.Get("method")
		)

		if method == "GET_LIST" {
			h.GetListCategory(w, r)
		} else if method == "GET" {
			h.GetByIdCategory(w, r)
		}
	case "PUT":
		h.UpdateCategory(w, r)
	case "DELETE":
		h.DeleteCategory(w, r)

	}
}

func (h *handler) CreateCategory(w http.ResponseWriter, r *http.Request) {

	var createCategory models.CreateCategory

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.handlerResponse(w, "error while read body: "+err.Error(), http.StatusBadRequest, nil)
		return
	}

	err = json.Unmarshal(body, &createCategory)
	if err != nil {
		h.handlerResponse(w, "error while unmarshal body: "+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	id, err := h.strg.Category().Create(&createCategory)
	if err != nil {
		h.handlerResponse(w, "error while storage category create:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	resp, err := h.strg.Category().GetByID(&models.CategoryPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(w, "error while storage category get by id:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(w, "Success", http.StatusOK, resp)
}

func (h *handler) GetByIdCategory(w http.ResponseWriter, r *http.Request) {

	var id string = r.URL.Query().Get("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(w, "invalid id uuid", http.StatusBadRequest, nil)
		return
	}

	resp, err := h.strg.Category().GetByID(&models.CategoryPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(w, "error while storage category get by id:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(w, "Success", http.StatusOK, resp)
}

func (h *handler) GetListCategory(w http.ResponseWriter, r *http.Request) {

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

	resp, err := h.strg.Category().GetList(&models.CategoryGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: search,
	})
	if err != nil {
		h.handlerResponse(w, "error while storage category get list:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(w, "Success", http.StatusOK, resp)
}

func (h *handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {

	var (
		resp = models.UpdateCategory{}
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

	cgryModel, err := h.strg.Category().Update(&resp)
	if err != nil {
		h.handlerResponse(w, "error while storage category update:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(w, "Succes", http.StatusOK, cgryModel)
}

func (h *handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	var (
		resp = models.CategoryPrimaryKey{}
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
	err = h.strg.Category().Delete(&resp)
	if err != nil {
		h.handlerResponse(w, "error while storage category delete:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(err.Error()))
}
