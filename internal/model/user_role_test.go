package model

import "testing"

func TestUserRole_Valid(t *testing.T) {
	tests := []struct {
		name string
		r    UserRole
		want bool
	}{
		{
			name: "Valid.employee",
			r:    UserRoleEmployee,
			want: true,
		},
		{
			name: "Valid.moderator",
			r:    UserRoleModerator,
			want: true,
		},
		{
			name: "invalid",
			r:    UserRole("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
