# lxc-go-http-api

An api to create and manage lxc through HTTP.

# Buildind

From the root source tree :

```
make
```

# Using

From the root source tree :

```
bin/lxc-go-http-api
```

After that, API listen on port 8000.

# Documentation

API documentation is in [OpenAPI 2.0](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md) format and generated with [go-swagger](https://goswagger.io/) command.

At the moment, there is two ways to access documentation :

* Read it to JSON format (**docs/swagger.json**) or copy-paste JSON content in [Swagger editor](https://editor.swagger.io)
* Run application from source tree (**bin/lxc-go-http-api**) and access documentation from this URL : http://server:8000/docs

*To access documentation from application we used a [Redoc middleware](https://github.com/go-openapi/runtime/blob/master/middleware/redoc.go).*


If you need to generate documentation, get [go-swagger](https://goswagger.io/) command:

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