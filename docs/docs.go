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
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/entries/create": {
            "post": {
                "description": "Create entry.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "entries"
                ],
                "summary": "Create entry.",
                "parameters": [
                    {
                        "description": "entry info",
                        "name": "entry",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/delivery.CreateEntryIn"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success create entry",
                        "schema": {
                            "$ref": "#/definitions/delivery.CreateEntryOut"
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "404": {
                        "description": "item is not found",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "422": {
                        "description": "unprocessable entity",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/me/entries": {
            "get": {
                "description": "Get my entries or get my entries for a day",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "entries"
                ],
                "summary": "Get my entries.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "day for events in YYYY-MM-DD format",
                        "name": "day",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success create entry",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/delivery.EntryOut"
                            }
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/me/projects": {
            "get": {
                "description": "Get my projects or get my projects for a day",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "projects"
                ],
                "summary": "Get my projects.",
                "responses": {
                    "200": {
                        "description": "success create project",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/delivery.ProjectOut"
                            }
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/projects/create": {
            "post": {
                "description": "Create project.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "projects"
                ],
                "summary": "Create project.",
                "parameters": [
                    {
                        "description": "project info",
                        "name": "project",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/delivery.CreateProjectIn"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success create project",
                        "schema": {
                            "$ref": "#/definitions/delivery.CreateProjectOut"
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "422": {
                        "description": "unprocessable entity",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "delivery.CreateEntryIn": {
            "type": "object",
            "required": [
                "project_id",
                "time_end",
                "time_start"
            ],
            "properties": {
                "name": {
                    "type": "string"
                },
                "project_id": {
                    "type": "integer"
                },
                "time_end": {
                    "type": "string"
                },
                "time_start": {
                    "type": "string"
                }
            }
        },
        "delivery.CreateEntryOut": {
            "type": "object",
            "required": [
                "id"
            ],
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "delivery.CreateProjectIn": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "delivery.CreateProjectOut": {
            "type": "object",
            "required": [
                "id"
            ],
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "delivery.EntryOut": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "project_id": {
                    "type": "integer"
                },
                "time_end": {
                    "type": "string"
                },
                "time_start": {
                    "type": "string"
                }
            }
        },
        "delivery.ProjectOut": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "echo.HTTPError": {
            "type": "object",
            "properties": {
                "message": {}
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
