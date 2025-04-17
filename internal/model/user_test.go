package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestParseUserID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		res, err := ParseUserID("95c386d5-d629-455d-994f-f64752bc3a2b")
		require.NoError(t, err)
		require.Equal(t, uuid.UUID{0x95, 0xc3, 0x86, 0xd5, 0xd6, 0x29, 0x45, 0x5d, 0x99, 0x4f, 0xf6, 0x47, 0x52, 0xbc, 0x3a, 0x2b}, res.UUID())
	})
	t.Run("error", func(t *testing.T) {
		_, err := ParseUserID("95c386d5-d629-455d-994f-f64752bc3a2b1")
		require.Error(t, err)
	})
}
