package repository

import (
	"database/sql"
	"strconv"

	"employee-service/models"
)

type EmployeeRepository struct {
	DB *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{DB: db}
}

func (r *EmployeeRepository) Create(employee *models.Employee) error {

	query := `
	INSERT INTO employees (
		first_name,
		last_name,
		date_of_birth,
		gender,
		email,
		phone_number,
		address,
		username,
		position,
		department,
		active
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,true)
	RETURNING id, created_at, updated_at
	`

	return r.DB.QueryRow(
		query,
		employee.FirstName,
		employee.LastName,
		employee.DateOfBirth,
		employee.Gender,
		employee.Email,
		employee.PhoneNumber,
		employee.Address,
		employee.Username,
		employee.Position,
		employee.Department,
	).Scan(&employee.ID, &employee.CreatedAt, &employee.UpdatedAt)
}

func (r *EmployeeRepository) GetAll(filter *models.EmployeeFilter) ([]models.Employee, error) {
	query := `
	SELECT id, first_name, last_name, date_of_birth, gender, email,
	       phone_number, address, username, position, department,
	       active, created_at, updated_at
	FROM employees
	WHERE 1=1
	`

	args := []any{}
	argPos := 1

	if filter.Email != "" {
		query += " AND email ILIKE $" + strconv.Itoa(argPos)
		args = append(args, "%"+filter.Email+"%")
		argPos++
	}

	if filter.FirstName != "" {
		query += " AND first_name ILIKE $" + strconv.Itoa(argPos)
		args = append(args, "%"+filter.FirstName+"%")
		argPos++
	}

	if filter.LastName != "" {
		query += " AND last_name ILIKE $" + strconv.Itoa(argPos)
		args = append(args, "%"+filter.LastName+"%")
		argPos++
	}

	if filter.Position != "" {
		query += " AND position ILIKE $" + strconv.Itoa(argPos)
		args = append(args, "%"+filter.Position+"%")
		argPos++
	}

	query += " ORDER BY id ASC"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []models.Employee

	for rows.Next() {
		var employee models.Employee

		err := rows.Scan(
			&employee.ID,
			&employee.FirstName,
			&employee.LastName,
			&employee.DateOfBirth,
			&employee.Gender,
			&employee.Email,
			&employee.PhoneNumber,
			&employee.Address,
			&employee.Username,
			&employee.Position,
			&employee.Department,
			&employee.Active,
			&employee.CreatedAt,
			&employee.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		employees = append(employees, employee)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *EmployeeRepository) GetByID(id int64) (*models.Employee, error) {
	query := `
	SELECT id, first_name, last_name, date_of_birth, gender, email,
	       phone_number, address, username, position, department,
	       active, created_at, updated_at
	FROM employees
	WHERE id = $1
	`

	var employee models.Employee

	err := r.DB.QueryRow(query, id).Scan(
		&employee.ID,
		&employee.FirstName,
		&employee.LastName,
		&employee.DateOfBirth,
		&employee.Gender,
		&employee.Email,
		&employee.PhoneNumber,
		&employee.Address,
		&employee.Username,
		&employee.Position,
		&employee.Department,
		&employee.Active,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &employee, nil
}

func (r *EmployeeRepository) Update(employee *models.Employee) error {
	query := `
	UPDATE employees
	SET first_name = $1,
	    last_name = $2,
	    date_of_birth = $3,
	    gender = $4,
	    email = $5,
	    phone_number = $6,
	    address = $7,
	    username = $8,
	    position = $9,
	    department = $10,
	    active = $11,
	    updated_at = NOW()
	WHERE id = $12
	RETURNING updated_at
	`

	return r.DB.QueryRow(
		query,
		employee.FirstName,
		employee.LastName,
		employee.DateOfBirth,
		employee.Gender,
		employee.Email,
		employee.PhoneNumber,
		employee.Address,
		employee.Username,
		employee.Position,
		employee.Department,
		employee.Active,
		employee.ID,
	).Scan(&employee.UpdatedAt)
}

func (r *EmployeeRepository) UpdateActiveStatus(id int64, active bool) error {
	query := `
	UPDATE employees
	SET active = $1,
	    updated_at = NOW()
	WHERE id = $2
	`
	_, err := r.DB.Exec(query, active, id)
	return err
}

func (r *EmployeeRepository) GetPermissions(employeeID int64) ([]string, error) {
	query := `
	SELECT permission
	FROM employee_permissions
	WHERE employee_id = $1
	ORDER BY permission ASC
	`

	rows, err := r.DB.Query(query, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string

	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *EmployeeRepository) ReplacePermissions(employeeID int64, permissions []string) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	deleteQuery := `DELETE FROM employee_permissions WHERE employee_id = $1`
	_, err = tx.Exec(deleteQuery, employeeID)
	if err != nil {
		return err
	}

	insertQuery := `
	INSERT INTO employee_permissions (employee_id, permission)
	VALUES ($1, $2)
	`

	for _, permission := range permissions {
		_, err = tx.Exec(insertQuery, employeeID, permission)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *EmployeeRepository) HasPermission(employeeID int64, permission string) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM employee_permissions
		WHERE employee_id = $1 AND permission = $2
	)
	`

	var exists bool
	err := r.DB.QueryRow(query, employeeID, permission).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}