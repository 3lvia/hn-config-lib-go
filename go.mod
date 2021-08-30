module github.com/3lvia/hn-config-lib-go

go 1.13

require (
	cloud.google.com/go/bigquery v1.3.0
	cloud.google.com/go/datastore v1.0.0
	github.com/MicahParks/keyfunc v0.7.0
	github.com/golang-jwt/jwt/v4 v4.0.0
	github.com/pkg/errors v0.8.1
	golang.org/x/net v0.0.0-20191126235420-ef20fe5d7933
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae // indirect
	google.golang.org/api v0.13.0
)

replace github.com/3lvia/hn-config-lib-go/elvid => ./elvid
