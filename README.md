# Xsolla Blog Post Web Service

## Test
Install db testing environment
```console
go get -u github.com/mattn/go-sqlite3    
```
Run the test
```console
make test
```

## Quick Run
Set some environment variables under [dev.sh](cmd/dev.sh)
```
export APP=blog_app
export URL_PREFIX=api/v1
export DB_USERNAME=root
export DB_PASSWORD=root
export DB_HOST=localhost
export DB_PORT=3307
export DB_DATABASE=local
export LOG_DIR=./log
```
Run app
```console
go mod tidy
make dev
```
## API Doc
[api__blogpost_v1.yml](api/openapi-spec/v3//api__blogpost_v1.yml)