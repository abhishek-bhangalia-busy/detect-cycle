package models

type EmployeeManager struct {
	EmployeeID uint64 `pg:",pk"`
	ManagerID uint64 `pg:",pk"`
}