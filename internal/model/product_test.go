package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestParseProductID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		res, err := ParseProductID("95c386d5-d629-455d-994f-f64752bc3a2b")
		require.NoError(t, err)
		require.Equal(t, uuid.UUID{0x95, 0xc3, 0x86, 0xd5, 0xd6, 0x29, 0x45, 0x5d, 0x99, 0x4f, 0xf6, 0x47, 0x52, 0xbc, 0x3a, 0x2b}, res.UUID())
	})
	t.Run("error", func(t *testing.T) {
		_, err := ParseProductID("95c386d5-d629-455d-994f-f64752bc3a2b1")
		require.Error(t, err)
	})
}

func TestParseProductCategory(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    ProductCategory
		wantErr bool
	}{
		{
			name:    "Valid.electronics",
			arg:     "электроника",
			want:    ProductCategoryElectronics,
			wantErr: false,
		},
		{
			name:    "Valid.clothes",
			arg:     "одежда",
			want:    ProductCategoryClothes,
			wantErr: false,
		},
		{
			name:    "Valid.shoes",
			arg:     "обувь",
			want:    ProductCategoryShoes,
			wantErr: false,
		},
		{
			name:    "invalid.0",
			arg:     "invalid",
			want:    ProductCategory(0),
			wantErr: true,
		},
		{
			name:    "empty",
			arg:     "",
			want:    ProductCategory(0),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParseProductCategory(tt.arg)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.want, res)
		})
	}
}
