package pb 

const (
  Swagger = `{
  "swagger": "2.0",
  "info": {
    "title": "fxgateway.proto",
    "version": "version not set"
  },
  "schemes": [
    "http",
    "https"
  ],
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "paths": {
    "/healthz": {
      "get": {
        "operationId": "HealthCheck",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/pbMessage"
            }
          }
        },
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/system/function-log/{FunctionName}": {
      "get": {
        "operationId": "GetLog",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/pbMessage"
            }
          }
        },
        "parameters": [
          {
            "name": "FunctionName",
            "in": "path",
            "required": true,
            "type": "string"
          }
        ],
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/system/function/{FunctionName}": {
      "get": {
        "operationId": "GetMeta",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/pbFunction"
            }
          }
        },
        "parameters": [
          {
            "name": "FunctionName",
            "in": "path",
            "required": true,
            "type": "string"
          }
        ],
        "tags": [
          "FxGateway"
        ]
      },
      "delete": {
        "operationId": "Delete",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/pbMessage"
            }
          }
        },
        "parameters": [
          {
            "name": "FunctionName",
            "in": "path",
            "required": true,
            "type": "string"
          }
        ],
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/system/functions": {
      "get": {
        "operationId": "List",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/pbFunctions"
            }
          }
        },
        "tags": [
          "FxGateway"
        ]
      },
      "post": {
        "operationId": "Deploy",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/pbMessage"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/pbCreateFunctionRequest"
            }
          }
        ],
        "tags": [
          "FxGateway"
        ]
      },
      "put": {
        "operationId": "Update",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/pbMessage"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/pbCreateFunctionRequest"
            }
          }
        ],
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/system/info": {
      "get": {
        "operationId": "Info",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/pbMessage"
            }
          }
        },
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/system/scale-function": {
      "put": {
        "operationId": "ReplicaUpdate",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/pbMessage"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/pbScaleServiceRequest"
            }
          }
        ],
        "tags": [
          "FxGateway"
        ]
      }
    }
  },
  "definitions": {
    "pbCreateFunctionRequest": {
      "type": "object",
      "properties": {
        "Service": {
          "type": "string"
        },
        "Image": {
          "type": "string"
        },
        "EnvVars": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        },
        "Labels": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        },
        "Annotations": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        },
        "Constraints": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "Secrets": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "RegistryAuth": {
          "type": "string"
        },
        "Limits": {
          "$ref": "#/definitions/pbFunctionResources"
        },
        "Requests": {
          "$ref": "#/definitions/pbFunctionResources"
        },
        "MinReplicas": {
          "type": "integer",
          "format": "int32"
        },
        "MaxReplicas": {
          "type": "integer",
          "format": "int32"
        }
      }
    },
    "pbFunction": {
      "type": "object",
      "properties": {
        "Name": {
          "type": "string"
        },
        "Image": {
          "type": "string"
        },
        "InvocationCount": {
          "type": "string",
          "format": "uint64"
        },
        "Replicas": {
          "type": "string",
          "format": "uint64"
        },
        "AvailableReplicas": {
          "type": "string",
          "format": "uint64"
        },
        "Annotations": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        },
        "Labels": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        }
      }
    },
    "pbFunctionResources": {
      "type": "object",
      "properties": {
        "Memory": {
          "type": "string"
        },
        "CPU": {
          "type": "string"
        },
        "GPU": {
          "type": "string"
        }
      }
    },
    "pbFunctions": {
      "type": "object",
      "properties": {
        "Functions": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/pbFunction"
          }
        }
      }
    },
    "pbMessage": {
      "type": "object",
      "properties": {
        "Msg": {
          "type": "string"
        }
      }
    },
    "pbScaleServiceRequest": {
      "type": "object",
      "properties": {
        "ServiceName": {
          "type": "string"
        },
        "Replicas": {
          "type": "string",
          "format": "uint64"
        }
      }
    }
  }
}
`
)
