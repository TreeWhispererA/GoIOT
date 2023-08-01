package middlewares

import (
	"encoding/json"
	http "github.com/helios/go-sdk/proxy-libs/helioshttp"
)

func AuthorizationResponse(msg string, writer http.ResponseWriter) {
	type errdata struct {
		Message string `json:"message"`
	}
	temp := &errdata{Message: msg}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(writer).Encode(temp)
}

func SuccessArrRespond(fields interface{}, modelType string, writer http.ResponseWriter) {
	_, err := json.Marshal(fields)
	type data struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}
	temp := &data{Data: fields, Message: "success"}
	if err != nil {
		ServerErrResponse(err.Error(), writer)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	switch modelType {
	// case "ObjectType":
	// 	temp.Data = fields.([]*models.ObjectType)
	// case "DeviceType":
	// 	temp.Data = fields.([]*models.DeviceType)
	default:
		// handle invalid model type
	}

	json.NewEncoder(writer).Encode(temp)
}

// SuccessMessageResponse -> success error messageformatter
func SuccessMessageResponse(msg string, writer http.ResponseWriter) {
	type errdata struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	temp := &errdata{Message: "success", Status: msg}

	//Send header, status code and output to writer
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(temp)
}

// ErrorResponse -> error formatter
func ErrorResponse(error string, writer http.ResponseWriter) {
	type errdata struct {
		Message string `json:"message"`
	}
	temp := &errdata{Message: error}

	//Send header, status code and output to writer
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(writer).Encode(temp)
}

// ServerErrResponse -> server error formatter
func ServerErrResponse(error string, writer http.ResponseWriter) {
	type servererrdata struct {
		Message string `json:"msg"`
	}
	temp := &servererrdata{Message: error}

	//Send header, status code and output to writer
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(writer).Encode(temp)
}

// ValidationResponse -> user input validation
func ValidationResponse(fields map[string][]string, writer http.ResponseWriter) {
	//Create a new map and fill it
	response := make(map[string]interface{})
	response["data"] = fields
	response["message"] = "validation error"

	//Send header, status code and output to writer
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(writer).Encode(response)
}
