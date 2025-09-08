package main

import (
	c_http "immodi/novel-site/internal/http"
	"net/http"
)

func main() {
	r := c_http.Router{}

	http.ListenAndServe("0.0.0.0:3000", r.NewRouter())
}
