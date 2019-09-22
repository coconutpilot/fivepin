package listener

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Listener is a container for daemon vars
type Listener struct {
	port   int
	Mux    *http.ServeMux
	Server *http.Server
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("viewHandler()")

	w.Header().Set("cache-control", "private, max-age=0, no-store")
	fmt.Fprintf(w, r.URL.String())
	time.Sleep(time.Second * 10)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("pingHandler()")

	w.Header().Set("cache-control", "private, max-age=0, no-store")

	switch r.Method {
	case "GET":
		log.Printf("pong: %s", r.URL.String())
		fmt.Fprintf(w, r.URL.String())

	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Body error: %s", err)
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
		log.Printf("pong: %s", body)
		fmt.Fprintf(w, "%s", body)

	case "PUT":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Body error: %s", err)
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
		log.Printf("pong: %s", body)
		fmt.Fprintf(w, "%s", body)

	default:
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}
	log.Println("pingHandler() exit")
}

// New - factory
func New(port *int) Listener {
	var d Listener
	d.port = *port
	srvAddr := fmt.Sprintf(":%d", d.port)

	d.Mux = http.NewServeMux()

	d.Mux.HandleFunc("/", viewHandler)
	d.Mux.HandleFunc("/ping/", pingHandler)

	d.Server = &http.Server{Addr: srvAddr, Handler: d.Mux}

	return d
}
