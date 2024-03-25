package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

type Employee struct {
	id       int64
	fullName string
}

type WorkTime struct {
	id         int64
	employeeId int64
	startDate  string
	startTime  string
	stopDate   string
	stopTime   string
}

var validationErrors = []string{}

func main() {
	server := os.Getenv("DB_SERVER")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE")
	driverName := "sqlserver"

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", server, user, password, port, database)
	conn, err := connectToDatabase(connString, driverName)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer conn.Close()

	employees, err := fetchEmployees(conn)

	if err != nil {
		log.Fatal(err)
		return
	}

	workTimes, err := fetchWorkTimes(conn)

	if err != nil {
		log.Fatal(err)
		return
	}

	for _, employee := range employees {
		if !validateEmployee(&employee) {
			continue
		}

		if err = saveEmployee(conn, &employee); err != nil {
			log.Fatal(err)
			return
		}
	}

	for _, workTime := range workTimes {
		if !validateWorkTime(&workTime) || !validateIntegrity(&workTime, &employees) {
			continue
		}

		if err = saveWorkTime(conn, &workTime); err != nil {
			log.Fatal(err)
			return
		}
	}

	for _, validationError := range validationErrors {
		if err := saveValidationError(conn, validationError); err != nil {
			log.Fatal(err)
			return
		}
	}
}

func connectToDatabase(connString string, driverName string) (*sql.DB, error) {
	conn, err := sql.Open(driverName, connString)

	if err != nil {
		return nil, err
	}

	err = conn.Ping()

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func fetchEmployees(conn *sql.DB) ([]Employee, error) {
	rows, err := conn.Query("SELECT * FROM dbo.temp_employees")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	employees := []Employee{}

	for rows.Next() {
		var employee Employee

		if err := rows.Scan(&employee.id, &employee.fullName); err != nil {
			return nil, err
		}

		employees = append(employees, employee)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}

func fetchWorkTimes(conn *sql.DB) ([]WorkTime, error) {
	rows, err := conn.Query("SELECT * FROM dbo.temp_worktimes")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	workTimes := []WorkTime{}

	for rows.Next() {
		var workTime WorkTime

		if err := rows.Scan(&workTime.id, &workTime.employeeId, &workTime.startDate, &workTime.startTime, &workTime.stopDate, &workTime.stopTime); err != nil {
			return nil, err
		}

		workTimes = append(workTimes, workTime)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return workTimes, nil
}

func validateEmployee(employee *Employee) bool {
	isValid := true

	if employee.id < 0 {
		logValidationError(fmt.Sprintf("Employee ID %d is invalid", employee.id))
		isValid = false
	}

	if len(employee.fullName) == 0 {
		logValidationError(fmt.Sprintf("Employee (id: %d); Employee full name is empty", employee.id))
		isValid = false
	}

	regex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !regex.MatchString(employee.fullName) {
		logValidationError(fmt.Sprintf("Employee (id: %d); Employee full name contains invalid characters", employee.id))
		isValid = false
	}

	return isValid
}

func validateWorkTime(workTime *WorkTime) bool {
	isValid := true

	if workTime.id < 0 {
		logValidationError(fmt.Sprintf("WorkTime ID %d is invalid", workTime.id))
		isValid = false
	}

	if workTime.employeeId < 0 {
		logValidationError(fmt.Sprintf("WorkTime (id: %d); Employee ID %d is invalid", workTime.id, workTime.employeeId))
		isValid = false
	}

	startDateTime, err := time.Parse("2.01.2006 15:04:05", fmt.Sprintf("%s %s", workTime.startDate, workTime.startTime))

	if err != nil {
		logValidationError(fmt.Sprintf("WorkTime (id: %d); Start date and time is invalid", workTime.id))
		isValid = false
	}

	stopDateTime, err := time.Parse("2.01.2006 15:04:05", fmt.Sprintf("%s %s", workTime.stopDate, workTime.stopTime))

	if err != nil {
		logValidationError(fmt.Sprintf("WorkTime (id: %d); Stop date and time is invalid", workTime.id))
		isValid = false
	}

	if startDateTime.After(stopDateTime) {
		logValidationError(fmt.Sprintf("WorkTime (id: %d); Start date and time is after stop date and time", workTime.id))
		isValid = false
	}

	return isValid
}

func validateIntegrity(workTime *WorkTime, employees *[]Employee) bool {
	isValid := false

	for _, employee := range *employees {
		if employee.id == workTime.employeeId {
			isValid = true
			break
		}
	}

	if !isValid {
		logValidationError(fmt.Sprintf("WorkTime (id: %d); Employee ID %d does not exist", workTime.id, workTime.employeeId))
	}

	return isValid
}

func logValidationError(message string) {
	validationErrors = append(validationErrors, message)
}

func saveEmployee(conn *sql.DB, employee *Employee) error {
	if _, err := conn.Exec("INSERT INTO dbo.employees (id, full_name) VALUES (?, ?)", employee.id, employee.fullName); err != nil {
		return err
	}

	return nil
}

func saveWorkTime(conn *sql.DB, workTime *WorkTime) error {
	if _, err := conn.Exec("INSERT INTO dbo.worktimes (id, employee_id, start_date, start_time, stop_date, stop_time) VALUES (?, ?, ?, ?, ?, ?)", workTime.id, workTime.employeeId, workTime.startDate, workTime.startTime, workTime.stopDate, workTime.stopTime); err != nil {
		return err
	}

	return nil
}

func saveValidationError(conn *sql.DB, message string) error {
	if _, err := conn.Exec("INSERT INTO dbo.validation_errors (message) VALUES (?)", message); err != nil {
		return err
	}

	return nil
}
