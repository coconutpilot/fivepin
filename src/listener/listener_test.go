package listener

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestViewHandlerGET(t *testing.T) {
	time := time.Now()
	uri := fmt.Sprintf("http://foo.example.com/nada/%v", time.String())
	r, _ := http.NewRequest("GET", uri, nil)
	w := httptest.NewRecorder()

	viewHandler(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
}

func TestPingHandlerGET(t *testing.T) {
	time := time.Now()
	uri := fmt.Sprintf("http://foo.example.com/ping?q=%s", time.String())
	r, _ := http.NewRequest("GET", uri, nil)
	w := httptest.NewRecorder()

	pingHandler(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}

	if w.Body.String() != uri {
		t.Errorf("expected %q but instead got %q", uri, w.Body.String())
	}
}

func TestPingHandlerPOST(t *testing.T) {
	time := fmt.Sprintf("%v", time.Now())
	uri := "http://foo.example.com/ping/"
	body := bytes.NewBufferString(time)
	r, _ := http.NewRequest("POST", uri, body)
	w := httptest.NewRecorder()

	pingHandler(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}

	if w.Body.String() != time {
		t.Errorf("expected %q but instead got %q", time, w.Body.String())
	}
}

type bodyFailer int

func (bodyFailer) Read(p []byte) (n int, err error) {
	return 0, errors.New("Testing error")
}

func TestPingHandlerPOSTError(t *testing.T) {
	uri := "http://foo.example.com/ping/"
	r, _ := http.NewRequest("POST", uri, bodyFailer(0))
	w := httptest.NewRecorder()

	pingHandler(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusInternalServerError)
	}
}

func TestPingHandlerPUT(t *testing.T) {
	time := fmt.Sprintf("%v", time.Now())
	uri := "http://foo.example.com/ping/"
	body := bytes.NewBufferString(time)
	r, _ := http.NewRequest("PUT", uri, body)
	w := httptest.NewRecorder()

	pingHandler(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}

	if w.Body.String() != time {
		t.Errorf("expected %q but instead got %q", time, w.Body.String())
	}
}

func TestPingHandlerPUTError(t *testing.T) {
	uri := "http://foo.example.com/ping/"
	r, _ := http.NewRequest("PUT", uri, bodyFailer(1))
	w := httptest.NewRecorder()

	pingHandler(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusInternalServerError)
	}
}

func TestPingHandlerBOGUS(t *testing.T) {
	uri := "http://foo.example.com/ping/"
	r, _ := http.NewRequest("BOGUS", uri, nil)
	w := httptest.NewRecorder()

	pingHandler(w, r)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestNew(t *testing.T) {
	port := 8080
	d := New(&port)

	if d.port != 8080 {
		t.Errorf("Wrong port: got %v want %v", d.port, port)
	}

}
