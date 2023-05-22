package lorm

import (
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	var s Server = &HTTPServer{}
	if err := http.ListenAndServe(":8887", s); err != nil {
		panic(err)
	}
}
