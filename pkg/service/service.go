package service

import (
	"encoding/json"
	"github.com/artrey/remux/pkg/remux"
	"github.com/artrey/simple-pages/pkg/service/dto"
	"github.com/artrey/simple-pages/pkg/storage"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Server struct {
	storage storage.Interface
	mux     *remux.ReMux
}

func New(storage storage.Interface) *Server {
	return &Server{
		storage: storage,
		mux:     remux.New(),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) Init(middlewares ...remux.Middleware) {
	re, err := regexp.Compile(`^/api/pages/(?P<id>\d+)$`)
	if err != nil {
		panic(err)
	}

	err = s.mux.RegisterPlain(http.MethodGet, "/api/pages", http.HandlerFunc(s.getPages), middlewares...)
	if err != nil {
		panic(err)
	}
	err = s.mux.RegisterRegex(http.MethodGet, re, http.HandlerFunc(s.getPage), middlewares...)
	if err != nil {
		panic(err)
	}
	err = s.mux.RegisterPlain(http.MethodPost, "/api/pages", http.HandlerFunc(s.createPage), middlewares...)
	if err != nil {
		panic(err)
	}
	err = s.mux.RegisterRegex(http.MethodPut, re, http.HandlerFunc(s.updatePage), middlewares...)
	if err != nil {
		panic(err)
	}
	err = s.mux.RegisterRegex(http.MethodDelete, re, http.HandlerFunc(s.deletePage), middlewares...)
	if err != nil {
		panic(err)
	}
}

func (s *Server) getPages(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	pages, err := s.storage.GetPages()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(dto.MakeUnknownError(err))
		return
	}

	dtos := make([]*dto.PageInfo, len(pages))
	for i, page := range pages {
		dtos[i] = dto.PageInfoFromModelPage(page)
	}
	err = encoder.Encode(dtos)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(dto.MakeUnknownError(err))
	}
}

func (s *Server) getPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	params, err := remux.PathParams(r.Context())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(dto.MakeError("no-params", err))
		return
	}

	id, err := strconv.ParseInt(params.Named["id"], 10, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(dto.MakeError("bad-param", err))
		return
	}

	page, err := s.storage.GetPageById(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		_ = encoder.Encode(dto.MakeError("not-found", err))
		return
	}

	err = encoder.Encode(dto.PageDetailFromModelPage(page))

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(dto.MakeUnknownError(err))
	}
}

func (s *Server) createPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	decoder := json.NewDecoder(r.Body)
	var data dto.PageUpdate
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(dto.MakeError("invalid-data", err))
		return
	}

	page, err := s.storage.CreatePage(data.Title, data.ImageUri, data.Text)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(dto.MakeUnknownError(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = encoder.Encode(dto.PageDetailFromModelPage(page))

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(dto.MakeUnknownError(err))
	}
}

func (s *Server) updatePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	params, err := remux.PathParams(r.Context())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(dto.MakeError("no-params", err))
		return
	}

	id, err := strconv.ParseInt(params.Named["id"], 10, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(dto.MakeError("bad-param", err))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var data dto.PageUpdate
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(dto.MakeError("invalid-data", err))
		return
	}

	page, err := s.storage.UpdatePageById(id, data.Title, data.ImageUri, data.Text)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		_ = encoder.Encode(dto.MakeError("not-found", err))
		return
	}

	err = encoder.Encode(dto.PageDetailFromModelPage(page))

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(dto.MakeUnknownError(err))
	}
}

func (s *Server) deletePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	params, err := remux.PathParams(r.Context())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(dto.MakeError("no-params", err))
		return
	}

	id, err := strconv.ParseInt(params.Named["id"], 10, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(dto.MakeError("bad-param", err))
		return
	}

	err = s.storage.DeletePageById(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		_ = encoder.Encode(dto.MakeError("not-found", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
