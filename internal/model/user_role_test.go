package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseUserRole(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    UserRole
		wantErr bool
	}{
		{
			name:    "Valid.moderator",
			arg:     "moderator",
			want:    UserRoleModerator,
			wantErr: false,
		},
		{
			name:    "Valid.employee",
			arg:     "employee",
			want:    UserRoleEmployee,
			wantErr: false,
		},
		{
			name:    "invalid.0",
			arg:     "invalid",
			want:    UserRole(0),
			wantErr: true,
		},
		{
			name:    "empty",
			arg:     "",
			want:    UserRole(0),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParseUserRole(tt.arg)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.want, res)
		})
	}
}
