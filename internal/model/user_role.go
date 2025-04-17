package model

import "errors"

type UserRole int16

const (
	UserRoleModerator UserRole = 1
	UserRoleEmployee  UserRole = 2
)

func (s UserRole) String() string {
	switch s {
	case UserRoleModerator:
		return "moderator"
	case UserRoleEmployee:
		return "employee"
	}
	return ""
}

func ParseUserRole(s string) (UserRole, error) {
	switch s {
	case "moderator":
		return UserRoleModerator, nil
	case "employee":
		return UserRoleEmployee, nil
	}
	return UserRole(0), errors.New("role not found")
}
