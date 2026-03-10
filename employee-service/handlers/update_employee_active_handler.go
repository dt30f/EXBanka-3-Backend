package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"employee-service/models"
	"employee-service/repository"
	"employee-service/services"
)

type UpdateEmployeeActiveHandler struct {
	Service *services.EmployeeService
}

func NewUpdateEmployeeActiveHandler(db *sql.DB) *UpdateEmployeeActiveHandler {
	repo := repository.NewEmployeeRepository(db)
	service := services.NewEmployeeService(repo, nil)

	return &UpdateEmployeeActiveHandler{
		Service: service,
	}
}

func (h *UpdateEmployeeActiveHandler) Handle(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/employees/")
	path = strings.TrimSuffix(path, "/active")
	path = strings.Trim(path, "/")

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "invalid employee id", http.StatusBadRequest)
		return
	}

	var req models.UpdateEmployeeActiveRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	employee, err := h.Service.UpdateEmployeeActiveStatus(id, req.Active)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(employee)
}