package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/shweta20039/shortme/short"
)

func Analytics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("read analytics request error. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg, _ := json.Marshal(errorResp{Msg: http.StatusText(http.StatusInternalServerError)})
		w.Write(errMsg)
		return
	}

	var analytics analytics
	err = json.Unmarshal(body, &analytics)
	if err != nil {
		log.Printf("parse analytics request error. %v", err)
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(errorResp{Msg: http.StatusText(http.StatusBadRequest)})
		w.Write(errMsg)
		return
	}

	count, err := short.Shorter.Analytics(analytics.ShortURL)
	if err != nil {
		return err
	}
	anaResp, _ := json.Marshal(ananalyticsResp{Count: count})
	w.Write(anaResp)
}
