#/bin/sh
# run from root of repo
go test -v -coverpkg=./... -coverprofile=profile.cov ./...
go tool cover -func profile.cov
