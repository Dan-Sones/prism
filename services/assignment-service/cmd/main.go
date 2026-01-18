package main

import (
	"assignment-service/internal/api/http"
	http2 "net/http"
)

func main() {

	router := http.NewRouter()
	http.RegisterRoutes(router, http.Controllers{})

	http2.ListenAndServe(":8081", router)
}
