// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "summary": "Authenticate a user",
                "operationId": "login",
                "parameters": [
                    {
                        "description": "Add login credentials",
                        "name": "login",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Login"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Login"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "summary": "Register a user",
                "operationId": "register",
                "parameters": [
                    {
                        "description": "Add account details",
                        "name": "register",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Account"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Account"
                        }
                    }
                }
            }
        },
        "/disposition": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "summary": "Get dispositions",
                "operationId": "get_dispositions",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Limit (max 100)",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "pass nothing, 'order' or 'quote'",
                        "name": "requestType",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Dispositions"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "summary": "Create a disposition, an order or a quote",
                "operationId": "create_disposition",
                "parameters": [
                    {
                        "description": "Add dispositions",
                        "name": "disposition",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Disposition"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/DispositionResponse"
                        }
                    }
                }
            }
        },
        "/disposition/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "summary": "Get disposition",
                "operationId": "get_disposition",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Disposition"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "Account": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "Disposition": {
            "type": "object",
            "required": [
                "createdBy",
                "customerId",
                "deliveryDate",
                "description",
                "requestType"
            ],
            "properties": {
                "createdBy": {
                    "type": "integer"
                },
                "currency": {
                    "type": "string"
                },
                "customerId": {
                    "type": "integer"
                },
                "deliveryDate": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "requestType": {
                    "type": "string",
                    "enum": [
                        "order",
                        "quote"
                    ]
                }
            }
        },
        "DispositionResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "Dispositions": {
            "type": "object",
            "properties": {
                "dispositions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Disposition"
                    }
                },
                "limit": {
                    "type": "integer",
                    "format": "int64"
                },
                "page": {
                    "type": "integer",
                    "format": "int64"
                },
                "total": {
                    "type": "integer",
                    "format": "int64"
                }
            }
        },
        "Login": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "SM Swagger API",
	Description:      "Swagger API for Business X.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
