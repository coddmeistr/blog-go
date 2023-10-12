package common

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxik12233/blog/types"
)

func ReturnAnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(types.Response{
		Error: []types.ErrorDetail{
			{
				ErrorType:    "ErrorTypeUnauthorized",
				ErrorMessage: "You are not allowed to access this path",
			},
		},
		Status:  http.StatusUnauthorized,
		Message: "Unauthorized access",
	})
}

func ReturnSimpleError(c *gin.Context, status int, err error) {
	//w.WriteHeader(status)
	//json.NewEncoder(w).Encode(types.SimpleError{
	//	Error: err.Error(),
	//	Code: 1,
	//})
}
