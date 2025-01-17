// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Amirali Kalhor",
            "url": "http://www.swagger.io/support",
            "email": "aakalhor2000@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/calculate": {
            "get": {
                "description": "Calculate the rebate amount for a transaction. Requires a valid transaction ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rebate Management"
                ],
                "summary": "Calculate rebate amount",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Transaction ID (UUID format)",
                        "name": "transaction_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Rebate amount successfully calculated",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "number"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid transaction or rebate calculation error",
                        "schema": {
                            "$ref": "#/definitions/http.CodeResponse"
                        }
                    },
                    "404": {
                        "description": "Transaction not found",
                        "schema": {
                            "$ref": "#/definitions/http.CodeResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/http.CodeResponse"
                        }
                    }
                }
            }
        },
        "/api/claim": {
            "post": {
                "description": "Submit a claim using a valid transaction ID. The claim must meet eligibility criteria.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rebate Management"
                ],
                "summary": "Submit a rebate claim",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Transaction ID (UUID format)",
                        "name": "transaction_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Claim successfully created",
                        "schema": {
                            "$ref": "#/definitions/http.RebateClaim"
                        }
                    },
                    "400": {
                        "description": "Invalid input or user not eligible",
                        "schema": {
                            "$ref": "#/definitions/http.CodeResponse"
                        }
                    },
                    "404": {
                        "description": "Transaction not found",
                        "schema": {
                            "$ref": "#/definitions/http.CodeResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/http.CodeResponse"
                        }
                    }
                }
            }
        },
        "/api/rebate": {
            "post": {
                "description": "Register a rebate program. Program names must be unique and include a percentage rebate.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Rebate Management"
                ],
                "summary": "Register a new rebate program",
                "parameters": [
                    {
                        "description": "Rebate Program Details",
                        "name": "rebate",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/http.RebateProgram"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Rebate program successfully created",
                        "schema": {
                            "$ref": "#/definitions/http.RebateProgram"
                        }
                    },
                    "400": {
                        "description": "Invalid input or duplicate program name",
                        "schema": {
                            "$ref": "#/definitions/http.CodeResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/http.CodeResponse"
                        }
                    }
                }
            }
        },
        "/api/reporting": {
            "post": {
                "description": "Generate a detailed report of claims within a specified date range.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Reporting"
                ],
                "summary": "Generate claims report",
                "parameters": [
                    {
                        "description": "Date Range (start_date and end_date in YYYY-MM-DD format)",
                        "name": "dateRange",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/http.DateRange"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Claims report successfully generated",
                        "schema": {
                            "$ref": "#/definitions/http.RebateClaimsReport"
                        }
                    },
                    "400": {
                        "description": "Invalid date range or input",
                        "schema": {
                            "$ref": "#/definitions/http.CodeResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/http.CodeResponse"
                        }
                    }
                }
            }
        },
        "/api/transaction": {
            "post": {
                "description": "Record a transaction for a rebate program. The transaction must reference an existing rebate program.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transaction Management"
                ],
                "summary": "Record a new transaction",
                "parameters": [
                    {
                        "description": "Transaction Details",
                        "name": "transaction",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/http.Transaction"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Transaction successfully recorded",
                        "schema": {
                            "$ref": "#/definitions/http.Transaction"
                        }
                    },
                    "400": {
                        "description": "Invalid input or unable to create transaction",
                        "schema": {
                            "$ref": "#/definitions/http.CodeResponse"
                        }
                    },
                    "404": {
                        "description": "Rebate program not found",
                        "schema": {
                            "$ref": "#/definitions/http.CodeResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/http.CodeResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "http.ClaimMetrics": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "count": {
                    "type": "integer"
                }
            }
        },
        "http.CodeResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "http.DateRange": {
            "type": "object",
            "required": [
                "end_date",
                "start_date"
            ],
            "properties": {
                "end_date": {
                    "type": "string"
                },
                "start_date": {
                    "type": "string"
                }
            }
        },
        "http.RebateClaim": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "date": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "transaction_id": {
                    "type": "string"
                }
            }
        },
        "http.RebateClaimsReport": {
            "type": "object",
            "properties": {
                "approved": {
                    "$ref": "#/definitions/http.ClaimMetrics"
                },
                "from": {
                    "type": "string"
                },
                "pending": {
                    "$ref": "#/definitions/http.ClaimMetrics"
                },
                "rejected": {
                    "$ref": "#/definitions/http.ClaimMetrics"
                },
                "to": {
                    "type": "string"
                },
                "total": {
                    "$ref": "#/definitions/http.ClaimMetrics"
                }
            }
        },
        "http.RebateProgram": {
            "type": "object",
            "required": [
                "eligibility_criteria",
                "end_date",
                "percentage",
                "program_name",
                "start_date"
            ],
            "properties": {
                "eligibility_criteria": {
                    "type": "boolean"
                },
                "end_date": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "percentage": {
                    "type": "number"
                },
                "program_name": {
                    "type": "string"
                },
                "start_date": {
                    "type": "string"
                }
            }
        },
        "http.Transaction": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "date": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "rebate_id": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Rebate Program Swagger API",
	Description:      "Comprehensive API documentation for the Rebate Program service. This service handles rebate creation, transaction management, and claims processing.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
