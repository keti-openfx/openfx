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
    "/api/createns/{NamespaceName}": {
      "post": {
        "operationId": "Create",
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
            "name": "NamespaceName",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/pbCreateNamespaceRequest"
            }
          }
        ],
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/api/delete/{FunctionName}": {
      "post": {
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
          },
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/pbDeleteFunctionRequest"
            }
          }
        ],
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/api/deploy": {
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
      }
    },
    "/api/exit": {
      "post": {
        "operationId": "ExitIDE",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/pbExitResponse"
            }
          }
        },
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/api/invoke/{Service}": {
      "post": {
        "operationId": "Invoke",
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
            "name": "Service",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/pbInvokeServiceRequest"
            }
          }
        ],
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/api/list": {
      "post": {
        "operationId": "List",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/pbFunctions"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/pbTokenRequest"
            }
          }
        ],
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/api/login": {
      "post": {
        "operationId": "Login",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/pbLoginResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/pbLoginRequest"
            }
          }
        ],
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/api/requestIDE-URL": {
      "post": {
        "operationId": "StartIDE",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/pbStartResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/pbStartRequest"
            }
          }
        ],
        "tags": [
          "FxGateway"
        ]
      }
    },
    "/api/update": {
      "post": {
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
        },
        "token": {
          "type": "string"
        }
      }
    },
    "pbCreateNamespaceRequest": {
      "type": "object",
      "properties": {
        "NamespaceName": {
          "type": "string"
        }
      }
    },
    "pbDeleteFunctionRequest": {
      "type": "object",
      "properties": {
        "FunctionName": {
          "type": "string"
        },
        "token": {
          "type": "string"
        }
      }
    },
    "pbExitResponse": {
      "type": "object",
      "properties": {
        "Msg": {
          "type": "string"
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
    "pbInvokeServiceRequest": {
      "type": "object",
      "properties": {
        "Service": {
          "type": "string"
        },
        "Input": {
          "type": "string",
          "format": "byte"
        },
        "token": {
          "type": "string"
        }
      }
    },
    "pbLoginRequest": {
      "type": "object",
      "properties": {
        "token": {
          "type": "string"
        },
        "Member": {
          "type": "string"
        }
      }
    },
    "pbLoginResponse": {
      "type": "object",
      "properties": {
        "Msg": {
          "type": "string"
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
        },
        "token": {
          "type": "string"
        }
      }
    },
    "pbStartRequest": {
      "type": "object",
      "properties": {
        "token": {
          "type": "string"
        }
      }
    },
    "pbStartResponse": {
      "type": "object",
      "properties": {
        "IDE": {
          "type": "string"
        }
      }
    },
    "pbTokenRequest": {
      "type": "object",
      "properties": {
        "token": {
          "type": "string"
        }
      }
    }
  }
}
`
)
