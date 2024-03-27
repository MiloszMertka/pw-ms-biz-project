USE msbiz
GO

IF OBJECT_ID(N'dbo.temp_worktimes', N'U') IS NOT NULL
BEGIN
	DROP TABLE dbo.temp_worktimes
END
GO

IF OBJECT_ID(N'dbo.temp_employees', N'U') IS NOT NULL
BEGIN
	DROP TABLE dbo.temp_employees
END
GO

IF OBJECT_ID(N'dbo.worktimes', N'U') IS NOT NULL
BEGIN
	DROP TABLE dbo.worktimes
END
GO

IF OBJECT_ID(N'dbo.employees', N'U') IS NOT NULL
BEGIN
	DROP TABLE dbo.employees
END
GO

IF OBJECT_ID(N'dbo.validation_errors', N'U') IS NOT NULL
BEGIN
	DROP TABLE dbo.validation_errors
END
GO

CREATE TABLE dbo.temp_employees (
	employee_id BIGINT PRIMARY KEY,
	full_name NVARCHAR(255),
)
GO

CREATE TABLE dbo.temp_worktimes (
	id BIGINT PRIMARY KEY,
	employee_id BIGINT,
	start_date NVARCHAR(255),
	start_time NVARCHAR(255),
	stop_date NVARCHAR(255),
	stop_time NVARCHAR(255),
)
GO

CREATE TABLE dbo.employees (
	employee_id BIGINT PRIMARY KEY,
	full_name NVARCHAR(255),
)
GO

CREATE TABLE dbo.worktimes (
	id BIGINT PRIMARY KEY,
	employee_id BIGINT,
	start_date DATE,
	start_time TIME,
	stop_date DATE,
	stop_time TIME,
	FOREIGN KEY (employee_id) REFERENCES dbo.employees(employee_id),
)
GO

CREATE TABLE dbo.validation_errors (
	id BIGINT PRIMARY KEY IDENTITY,
	message NVARCHAR(255),
)
GO
