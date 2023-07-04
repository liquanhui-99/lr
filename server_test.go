package lorm

import "testing"

func TestHTTPServer(t *testing.T) {
	var h Server = &HTTPServer{}
	if err := h.Start(":8080"); err != nil {
		panic(err)
	}
}
