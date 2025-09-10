package main

import (
	"immodi/novel-site/internal/config"
	c_http "immodi/novel-site/internal/http"
	"net/http"
)

func main() {
	r := c_http.Router{}
	http.ListenAndServe("0.0.0.0:"+config.Port, r.NewRouter())
}
