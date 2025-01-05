module github.com/vaziolabs/LumberJack

replace dashboard => ./dashboard

go 1.23.2

require (
	dashboard v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/gorilla/mux v1.8.1
)

require golang.org/x/crypto v0.31.0
