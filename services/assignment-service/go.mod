module assignment-service

go 1.25

require (
	github.com/go-chi/chi/v5 v5.2.4
	github.com/go-chi/cors v1.2.2
	github.com/joho/godotenv v1.5.1
	github.com/Dan-Sones/prismdbmodels v0.0.0
)

replace github.com/Dan-Sones/prismdbmodels => ../prismdbmodels
