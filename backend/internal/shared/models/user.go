package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleSuperAdmin     Role = "super_admin"
	RoleHospitalAdmin  Role = "hospital_admin"
	RoleDepartmentHead Role = "department_head"
	RoleChargeNurse    Role = "charge_nurse"
	RoleStaff          Role = "staff"
	RoleEngineer       Role = "engineer"
)

type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	TenantID     uuid.UUID `db:"tenant_id" json:"tenant_id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     string    `db:"last_name" json:"last_name"`
	Role         Role      `db:"role" json:"role"`
	DepartmentID *uuid.UUID `db:"department_id" json:"department_id"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type Tenant struct {
	ID             uuid.UUID `db:"id" json:"id"`
	Name           string    `db:"name" json:"name"`
	Slug           string    `db:"slug" json:"slug"`
	Plan           string    `db:"plan" json:"plan"`
	MaxDevices     int       `db:"max_devices" json:"max_devices"`
	MaxDepartments int       `db:"max_departments" json:"max_departments"`
	IsActive       bool      `db:"is_active" json:"is_active"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}
