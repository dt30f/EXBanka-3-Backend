package middleware

import (
	"net/http"

	"employee-service/services"
)

func AdminOnlyMiddleware(employeeService *services.EmployeeService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			employeeID, ok := r.Context().Value(ContextEmployeeIDKey).(int64)
			if !ok || employeeID <= 0 {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			hasAdminPermission, err := employeeService.HasPermission(employeeID, "admin")
			if err != nil {
				http.Error(w, "failed to verify permissions", http.StatusInternalServerError)
				return
			}

			if !hasAdminPermission {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}