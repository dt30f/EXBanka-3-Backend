package services

import (
	"errors"
	"strings"

	"employee-service/models"
	"employee-service/repository"
)

type EmployeeService struct {
	EmployeeRepo *repository.EmployeeRepository
	AuthClient   *AuthClient
}

func NewEmployeeService(repo *repository.EmployeeRepository, authClient *AuthClient) *EmployeeService {
	return &EmployeeService{
		EmployeeRepo: repo,
		AuthClient:   authClient,
	}
}

func (s *EmployeeService) CreateEmployee(req *models.CreateEmployeeRequest) (*models.Employee, *models.CreateCredentialResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" {
		return nil, nil, errors.New("email is required")
	}

	if req.FirstName == "" || req.LastName == "" {
		return nil, nil, errors.New("first name and last name are required")
	}

	employee := &models.Employee{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DateOfBirth: req.DateOfBirth,
		Gender:      req.Gender,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Username:    req.Username,
		Position:    req.Position,
		Department:  req.Department,
		Active:      true,
	}

	err := s.EmployeeRepo.Create(employee)
	if err != nil {
		return nil, nil, err
	}

	credentialResponse, err := s.AuthClient.CreateCredential(employee.ID, employee.Email, false)
	if err != nil {
		return nil, nil, err
	}

	return employee, credentialResponse, nil
}

func (s *EmployeeService) GetAllEmployees(filter *models.EmployeeFilter) ([]models.Employee, error) {
	return s.EmployeeRepo.GetAll(filter)
}

func (s *EmployeeService) GetEmployeeByID(id int64) (*models.Employee, error) {
	if id <= 0 {
		return nil, errors.New("invalid employee id")
	}

	return s.EmployeeRepo.GetByID(id)
}

func (s *EmployeeService) UpdateEmployee(id int64, req *models.UpdateEmployeeRequest) (*models.Employee, error) {
	if id <= 0 {
		return nil, errors.New("invalid employee id")
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" {
		return nil, errors.New("email is required")
	}

	if req.FirstName == "" || req.LastName == "" {
		return nil, errors.New("first name and last name are required")
	}

	existingEmployee, err := s.EmployeeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	existingEmployee.FirstName = req.FirstName
	existingEmployee.LastName = req.LastName
	existingEmployee.DateOfBirth = req.DateOfBirth
	existingEmployee.Gender = req.Gender
	existingEmployee.Email = req.Email
	existingEmployee.PhoneNumber = req.PhoneNumber
	existingEmployee.Address = req.Address
	existingEmployee.Username = req.Username
	existingEmployee.Position = req.Position
	existingEmployee.Department = req.Department
	existingEmployee.Active = req.Active

	err = s.EmployeeRepo.Update(existingEmployee)
	if err != nil {
		return nil, err
	}

	return existingEmployee, nil
}

func (s *EmployeeService) UpdateEmployeeActiveStatus(id int64, active bool) (*models.Employee, error) {
	if id <= 0 {
		return nil, errors.New("invalid employee id")
	}

	employee, err := s.EmployeeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	err = s.EmployeeRepo.UpdateActiveStatus(id, active)
	if err != nil {
		return nil, err
	}

	employee.Active = active
	return employee, nil
}

func (s *EmployeeService) GetEmployeePermissions(id int64) ([]string, error) {
	if id <= 0 {
		return nil, errors.New("invalid employee id")
	}

	_, err := s.EmployeeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.EmployeeRepo.GetPermissions(id)
}

func (s *EmployeeService) UpdateEmployeePermissions(id int64, permissions []string) ([]string, error) {
	if id <= 0 {
		return nil, errors.New("invalid employee id")
	}

	_, err := s.EmployeeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	cleaned := make([]string, 0, len(permissions))
	seen := map[string]bool{}

	for _, p := range permissions {
		p = strings.TrimSpace(strings.ToLower(p))
		if p == "" {
			continue
		}
		if !seen[p] {
			seen[p] = true
			cleaned = append(cleaned, p)
		}
	}

	err = s.EmployeeRepo.ReplacePermissions(id, cleaned)
	if err != nil {
		return nil, err
	}

	return cleaned, nil
}

func (s *EmployeeService) HasPermission(employeeID int64, permission string) (bool, error) {
	if employeeID <= 0 {
		return false, errors.New("invalid employee id")
	}

	return s.EmployeeRepo.HasPermission(employeeID, permission)
}