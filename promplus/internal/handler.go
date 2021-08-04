package internal

import (
	"net/http"
	"strconv"
	"time"
)

func Handle(rw http.ResponseWriter, r *http.Request) {
	seconds, _ := strconv.Atoi(r.URL.Query().Get("duration"))
	status, _ := strconv.Atoi(r.URL.Query().Get("status"))

	if seconds > 0 {
		time.Sleep(time.Duration(seconds) * time.Second)
	}

	if status >= http.StatusOK {
		rw.WriteHeader(status)
	} else {
		rw.WriteHeader(http.StatusOK)
	}
}
