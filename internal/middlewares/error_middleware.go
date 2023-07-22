package middlewares

import (
	"fmt"
	"home-server/pkg/utils"
	"net/http"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("Global error handler:", err)
				if errorMap, ok := err.(utils.Errors); ok {
					errorMap.ToResponse(w)
				} else {
					// Handle the panic and build the error response
					errorMsg := "Internal Server Error: " + fmt.Sprintf("%v", err)

					// Set the appropriate headers and status code
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(errorMsg))
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
