package middlewares

import (
	"fmt"
	"net/http"

	"github.com/Brix101/network-file-manager/internals/utils"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("Error handler:", err)
				if errors, ok := err.(utils.Errors); ok {
					errors.Raise()
				} else {
					errorMsg := "Internal Server Error: " + fmt.Sprintf("%v", err)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(errorMsg))
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
