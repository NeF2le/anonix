package domain

type RoleName string

const (
	RoleAdmin      RoleName = "admin"
	RoleSpecialist RoleName = "specialist"
	RoleAuditor    RoleName = "auditor"
)
