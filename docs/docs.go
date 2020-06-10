// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/containers": {
            "get": {
                "description": "Return list of containers",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "containers"
                ],
                "summary": "Get containers list",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Containers"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.HTTPClientResp"
                        }
                    }
                }
            }
        },
        "/create": {
            "post": {
                "description": "Create a new container",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "container"
                ],
                "summary": "Create a new container",
                "parameters": [
                    {
                        "description": "Creation parameters",
                        "name": "options",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/main.LxcTemplate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.HTTPClientResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/main.HTTPClientResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.HTTPClientResp"
                        }
                    }
                }
            }
        },
        "/destroy/{container}": {
            "delete": {
                "description": "Destroy a container",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "container"
                ],
                "summary": "Destroy a container",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Container name",
                        "name": "container",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Destroy container even if it running",
                        "name": "force",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/main.APIDestroyOptions"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.HTTPClientResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/main.HTTPClientResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.HTTPClientResp"
                        }
                    }
                }
            }
        },
        "/version": {
            "get": {
                "description": "Return LXC version in used",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "general"
                ],
                "summary": "Get LXC version",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Version"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.HTTPClientResp"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.APIDestroyOptions": {
            "type": "object",
            "properties": {
                "force": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "main.Containers": {
            "type": "object",
            "properties": {
                "containers": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "main.HTTPClientResp": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "a successfully message"
                },
                "status": {
                    "type": "string",
                    "example": "success"
                }
            }
        },
        "main.LxcTemplate": {
            "$ref": "#/definitions/lxc.TemplateOptions"
        },
        "main.Version": {
            "type": "object",
            "properties": {
                "version": {
                    "type": "string",
                    "example": "4.0"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.1",
	Host:        "server.clerc.im:8000",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "LXC HTTP API",
	Description: "An api to create and manage lxc thourgh HTTP",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
