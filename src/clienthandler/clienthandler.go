package clienthandler

import (
	"fmt"
	"listener"
	"log"
	"net/http"
	"time"
)

// AddHandlers - client specific handlers
func AddHandlers(l *listener.Listener) {
	l.Mux.HandleFunc("/foo/", fooHandler)

}

func fooHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("fooHandler()")

	w.Header().Set("cache-control", "private, max-age=0, no-store")
	fmt.Fprintf(w, r.URL.String())
	time.Sleep(time.Second * 10)
}
