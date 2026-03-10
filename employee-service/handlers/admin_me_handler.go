package handlers

import (
	"encoding/json"
	"net/http"

	"employee-service/middleware"
)

func AdminMeHandler(w http.ResponseWriter, r *http.Request) {
	employeeID, _ := r.Context().Value(middleware.ContextEmployeeIDKey).(int64)
	email, _ := r.Context().Value(middleware.ContextEmailKey).(string)

	response := map[string]any{
		"message":     "admin access granted",
		"employee_id": employeeID,
		"email":       email,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}