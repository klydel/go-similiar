// Copyright 2013 Adam Peck

package api

import (
	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
	"encoding/json"
	"fmt"
	"go-similiar/img"
	"go-similiar/util"
	"net/http"
)

func init() {
	http.HandleFunc("/api/img", Handler)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		Delete(w, r)
	case "GET":
		Get(w, r)
	case "HEAD":
		util.Head(w, r, http.HandlerFunc(Get))
	case "OPTIONS":
		w.Header().Set("Allow", "DELETE, GET, HEAD, OPTIONS, POST, PUT, TRACE")
	case "POST":
		Post(w, r)
	case "PUT":
		Put(w, r)
	case "TRACE":
		util.Trace(w, r)
	default:
		util.NotImplemented(w)
	}
}

type ImgDAO struct {
	Context appengine.Context
}

func (dao ImgDAO) Delete(url string) error {
	return datastore.Delete(dao.Context, datastore.NewKey(dao.Context, "Img", url, 0, nil))
}

func (dao ImgDAO) Exists(url string) (bool, error) {
	var i img.Img
	return dao.Get(url, &i)
}

func (dao ImgDAO) Get(url string, i *img.Img) (bool, error) {
	err := datastore.Get(dao.Context, datastore.NewKey(dao.Context, "Img", url, 0, nil), i)
	if err == datastore.ErrNoSuchEntity {
		return false, nil
	}
	return err == nil, err
}

func (dao ImgDAO) Put(url string, i img.Img) error {
	_, err := datastore.Put(dao.Context, datastore.NewKey(dao.Context, "Img", url, 0, nil), &i)
	return err
}

func Get(w http.ResponseWriter, r *http.Request) {
	var i img.Img
	dao := ImgDAO{Context: appengine.NewContext(r)}
	ok, err := dao.Get(r.FormValue("url"), &i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	b, err := json.Marshal(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
	w.WriteHeader(http.StatusOK)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	dao := ImgDAO{Context: appengine.NewContext(r)}
	if err := dao.Delete(r.FormValue("url")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func Post(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	url := r.FormValue("url")
	i, err := img.Fetch(url, urlfetch.Client(c))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dao := ImgDAO{Context: c}
	if err := dao.Put(url, i); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", fmt.Sprint("?url=", url))
	w.WriteHeader(http.StatusNoContent)
}

func Put(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	url := r.FormValue("url")
	code := new(int)
	if err := datastore.RunInTransaction(c, func(c appengine.Context) error {
			dao := ImgDAO{Context: c}
			ok, err := dao.Exists(url)
			if err != nil {
				return err
			}

			if ok {
				*code = http.StatusNotModified
				return nil
			}
			*code = http.StatusNoContent

			i, err := img.Fetch(url, urlfetch.Client(c))
			if err != nil {
				return err
			}
			return dao.Put(url, i)
	}, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", fmt.Sprint("?url=", url))
	w.WriteHeader(*code)
}
