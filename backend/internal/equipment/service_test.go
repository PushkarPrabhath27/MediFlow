package equipment

import (
	"testing"

	"github.com/mediflow/backend/internal/shared/models"
	"github.com/stretchr/testify/assert"
)

func TestIsValidTransition(t *testing.T) {
	tests := []struct {
		name     string
		from     models.EquipmentStatus
		to       models.EquipmentStatus
		expected bool
	}{
		{"Available to In Use", models.StatusAvailable, models.StatusInUse, true},
		{"Available to Missing", models.StatusAvailable, models.StatusMissing, true},
		{"In Use to Available", models.StatusInUse, models.StatusAvailable, true},
		{"In Use to Reserved", models.StatusInUse, models.StatusReserved, false},
		{"Reserved to In Use", models.StatusReserved, models.StatusInUse, true},
		{"Reserved to Maintenance", models.StatusReserved, models.StatusInMaintenance, false},
		{"Maintenance to Available", models.StatusInMaintenance, models.StatusAvailable, true},
		{"Maintenance to Decommissioned", models.StatusInMaintenance, models.StatusDecommissioned, true},
		{"Transit to Available", models.StatusInTransit, models.StatusAvailable, true},
		{"Transit to In Use", models.StatusInTransit, models.StatusInUse, true},
		{"Missing to Available", models.StatusMissing, models.StatusAvailable, true},
		{"Decommissioned to Available", models.StatusDecommissioned, models.StatusAvailable, false},
		{"Same status", models.StatusAvailable, models.StatusAvailable, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidTransition(tt.from, tt.to)
			assert.Equal(t, tt.expected, result)
		})
	}
}
