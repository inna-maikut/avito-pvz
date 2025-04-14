package model

type UserRole string

const (
	UserRoleModerator UserRole = "moderator"
	UserRoleEmployee  UserRole = "employee"
)

func (r UserRole) Valid() bool {
	if r == UserRoleModerator || r == UserRoleEmployee {
		return true
	}

	return false
}
