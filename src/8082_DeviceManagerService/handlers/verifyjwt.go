package middlewares

import (
	"encoding/json"
	"fmt"

	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	http "github.com/helios/go-sdk/proxy-libs/helioshttp"
)

var mySigningKey = []byte(DotEnvVariable("JWT_SECRET"))

type Response struct {
	Message string `json:"message"`
}

// IsAuthorized -> verify jwt header
func IsAuthorized(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		if tokenString = c.Request.Header.Get("Token"); tokenString != "" {

		} else {
			AuthorizationResponse("Not Authorized", c.Writer)
			return
		}

		url := "http://localhost/api/v1/userservice/getinfo"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			AuthorizationResponse("Not Authorized.", c.Writer)
			return
		}
		req.Header.Set("Token", tokenString)
		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			AuthorizationResponse("Not Authorized.", c.Writer)
			return
		}
		defer resp.Body.Close()

		var response Response
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			fmt.Println("Error decoding response:", err)
			AuthorizationResponse("Not Authorized.", c.Writer)
			return
		}

		if response.Message != "success" {
			fmt.Println("Message:", response.Message)
			AuthorizationResponse(response.Message, c.Writer)
			return
		}
		next(c)
	}

}
