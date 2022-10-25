// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/upload": {
            "post": {
                "description": "upload csv file",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "upload csv file",
                "operationId": "uploader",
                "parameters": [
                    {
                        "type": "file",
                        "description": "query test file",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "The request was validated and has been processed successfully (sync)",
                        "schema": {
                            "$ref": "#/definitions/model.JSONSuccessResult"
                        }
                    },
                    "404": {
                        "description": "The payload was rejected as invalid",
                        "schema": {
                            "$ref": "#/definitions/model.JSONFailureResult"
                        }
                    },
                    "500": {
                        "description": "An internal error has occurred, most likely due to an uncaught exception",
                        "schema": {
                            "$ref": "#/definitions/model.JSONFailureResult"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.JSONFailureResult": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "data": {
                    "type": "object"
                },
                "error": {
                    "type": "string",
                    "example": "There was an error processing the request"
                },
                "id": {
                    "type": "string",
                    "example": "705e4dcb-3ecd-24f3-3a35-3e926e4bded5"
                },
                "stacktrace": {
                    "type": "string"
                }
            }
        },
        "model.JSONSuccessResult": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 200
                },
                "correlation_id": {
                    "type": "string",
                    "example": "705e4dcb-3ecd-24f3-3a35-3e926e4bded5"
                },
                "data": {
                    "type": "object"
                },
                "id": {
                    "type": "string",
                    "example": "123-456-789-abc-def"
                },
                "message": {
                    "type": "string",
                    "example": "Success"
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
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
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
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
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