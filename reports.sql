CREATE OR ALTER PROCEDURE AbsentEmployeesReport
    @Month INT,
    @Year INT
AS
BEGIN
	SELECT
		e.employee_id AS ID,
		e.full_name AS Osoba
	FROM
		dbo.employees e
	EXCEPT
	SELECT
		e.employee_id AS ID,
		e.full_name AS Osoba
	FROM dbo.employees e
		INNER JOIN dbo.worktimes w ON e.employee_id = w.employee_id
	WHERE
		DATEPART(MONTH, w.start_date) = @Month
		AND DATEPART(YEAR, w.start_date) = @Year
	GROUP BY
		e.employee_id,
		e.full_name
END
GO

CREATE OR ALTER PROCEDURE WorkHoursReport
    @Month INT,
    @Year INT
AS
BEGIN
    SELECT 
        e.employee_id AS ID,
        e.full_name AS Osoba,
        SUM(DATEDIFF(SECOND, CONVERT(DATETIME, w.start_date) + CONVERT(DATETIME, w.start_time), CONVERT(DATETIME, w.stop_date) + CONVERT(DATETIME, w.stop_time))) / 3600.0 AS Godziny
    FROM 
        dbo.employees e
    INNER JOIN 
        dbo.worktimes w ON e.employee_id = w.employee_id
    WHERE 
        DATEPART(MONTH, w.start_date) = @Month
        AND DATEPART(YEAR, w.start_date) = @Year
    GROUP BY 
        e.employee_id,
        e.full_name
    ORDER BY
        Godziny DESC
END
GO

CREATE OR ALTER PROCEDURE OvertimeReport
    @Month INT,
    @Year INT
AS
BEGIN
    SELECT 
        e.employee_id AS ID,
        e.full_name AS Osoba,
        SUM(
            CASE 
                WHEN DATEDIFF(SECOND, CONVERT(DATETIME, w.start_date) + CONVERT(DATETIME, w.start_time), CONVERT(DATETIME, w.stop_date) + CONVERT(DATETIME, w.stop_time)) > 8 * 3600 THEN 
                    DATEDIFF(SECOND, CONVERT(DATETIME, w.start_date) + CONVERT(DATETIME, w.start_time), CONVERT(DATETIME, w.stop_date) + CONVERT(DATETIME, w.stop_time)) / 3600.0 - 8
                ELSE 
                    0.0
            END
        ) AS Nadgodziny
    FROM 
        dbo.employees e
    INNER JOIN 
        dbo.worktimes w ON e.employee_id = w.employee_id
    WHERE 
        DATEPART(MONTH, w.start_date) = @Month
        AND DATEPART(YEAR, w.start_date) = @Year
    GROUP BY 
        e.employee_id,
        e.full_name
    ORDER BY
        Nadgodziny DESC
END
GO

EXEC AbsentEmployeesReport @Month = 1, @Year = 2023
GO

EXEC WorkHoursReport @Month = 1, @Year = 2023
GO

EXEC OvertimeReport @Month = 1, @Year = 2023
GO
