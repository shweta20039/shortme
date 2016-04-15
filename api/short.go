package api

import (
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/andyxning/shortme/short"
	"github.com/andyxning/shortme/conf"

	"github.com/gorilla/mux"
)

func redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortededURL := vars["shortenedURL"]

	longURL, err := short.Shorter.Expand(shortededURL)
	if err != nil {
		log.Printf("redirect short url error. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	} else {
		w.Header().Set("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func shortURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read short request error. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg, _ := json.Marshal(Err{Msg: http.StatusText(http.StatusInternalServerError)})
		w.Write(errMsg)
		return
	}

	var shortReq ShortReq
	err = json.Unmarshal(body, &shortReq)
	if err != nil {
		log.Printf("parse short request error. %v", err)
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(Err{Msg: http.StatusText(http.StatusBadRequest)})
		w.Write(errMsg)
		return
	}

	var shortenedURL string
	shortenedURL, err = short.Shorter.Short(shortReq.LongURL)
	if err != nil {
		log.Printf("short url error. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg, _ := json.Marshal(Err{Msg: http.StatusText(http.StatusInternalServerError)})
		w.Write(errMsg)
		return
	} else {
		shortResp, _ := json.Marshal(ShortResp{ShortURL: shortenedURL})
		w.Write(shortResp)
	}
}

func expandURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read expand request error. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg, _ := json.Marshal(Err{Msg: http.StatusText(http.StatusInternalServerError)})
		w.Write(errMsg)
		return
	}

	var expandReq ExpandReq
	err = json.Unmarshal(body, &expandReq)
	if err != nil {
		log.Printf("parse expand request error. %v", err)
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(Err{Msg: http.StatusText(http.StatusBadRequest)})
		w.Write(errMsg)
		return
	}

	var expandedURL string
	expandedURL, err = short.Shorter.Expand(expandReq.ShortURL)
	if err != nil {
		log.Printf("expand url error. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg, _ := json.Marshal(Err{Msg: http.StatusText(http.StatusInternalServerError)})
		w.Write(errMsg)
		return
	} else {
		expandResp, _ := json.Marshal(ExpandResp{LongURL: expandedURL})
		w.Write(expandResp)
	}
}


func Start() {
	log.Println("api starts")
	r := mux.NewRouter()
	r.HandleFunc("/version", version).Methods(http.MethodGet)
	r.HandleFunc("/health", healthCheck).Methods(http.MethodGet)
	r.HandleFunc("/short", shortURL).Methods(http.MethodPost).HeadersRegexp("Content-Type",	"application/json")
	r.HandleFunc("/expand", expandURL).Methods(http.MethodPost).HeadersRegexp("Content-Type", "application/json")
	r.HandleFunc("/{shortenedURL:[a-zA-Z0-9]{1,11}}", redirect).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(conf.Conf.Http.Listen, r))
}
