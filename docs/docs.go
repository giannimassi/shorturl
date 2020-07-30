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
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Gianni Massi",
            "url": "http://www.shorturl.com/support",
            "email": "support@shorturl.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api": {
            "get": {
                "description": "Returns information about the short url association stored for the provided key",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Return short URL info",
                "parameters": [
                    {
                        "description": "Key for which the request is made",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/routes.infoRequestPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.infoResponsePayload"
                        }
                    },
                    "400": {
                        "description": "Payload cannot be decoded"
                    },
                    "404": {
                        "description": "Key not found"
                    },
                    "500": {
                        "description": "The server has encountered an unknown error"
                    }
                }
            },
            "put": {
                "description": "Adds a new key-url association",
                "consumes": [
                    "application/json"
                ],
                "summary": "Add short url",
                "parameters": [
                    {
                        "description": "Key-url association to add",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/routes.addURLRequestPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Key-url association added"
                    },
                    "400": {
                        "description": "Payload cannot be decoded"
                    },
                    "409": {
                        "description": "A key-url association already exists for the provided key"
                    },
                    "422": {
                        "description": "URL in the payload is malformed"
                    },
                    "500": {
                        "description": "The server has encountered an unknown error"
                    }
                }
            },
            "delete": {
                "description": "Deletes a key-url association",
                "consumes": [
                    "application/json"
                ],
                "summary": "Delete short url",
                "parameters": [
                    {
                        "description": "Key-url association to delete",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/routes.deleteURLRequestPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Key-url association deleted"
                    },
                    "400": {
                        "description": "Payload cannot be decoded"
                    },
                    "404": {
                        "description": "Key-url association not found for key"
                    },
                    "500": {
                        "description": "The server has encountered an unknown error"
                    }
                }
            }
        }
    },
    "definitions": {
        "routes.addURLRequestPayload": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "routes.deleteURLRequestPayload": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                }
            }
        },
        "routes.infoRequestPayload": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                }
            }
        },
        "routes.infoResponsePayload": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
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
	Host:        "localhost:8080",
	BasePath:    "/api",
	Schemes:     []string{},
	Title:       "Shorturl API",
	Description: "This is an url shortening service",
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
