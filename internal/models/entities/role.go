package models

type Role struct {
	Name           string    `db:"name" json:"name"`
	PermissionCodes []string `db:"permission_codes" json:"permission_codes"`
	IsSystemRole   bool      `db:"is_system_role" json:"is_system_role"`
}