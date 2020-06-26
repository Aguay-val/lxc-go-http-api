# lxc-go-http-api

An api to create and manage lxc thourgh HTTP.

# Buildind

From the root source tree :

```
make build
```

# Using

```
bin/lxc-go-http-api
```

# Documentation

API documentation is generated with go-swagger 

If you need to generate documenation you need [go-swagger](https://goswagger.io/) command. Run :

```
go get -u github.com/go-swagger/go-swagger/cmd/swagger@v0.24.0
```

From the root source tree, run :

```
make swagger-generate
make swagger-validate
```

Or:

```
make docs
```
