package request

import (
	"context"
	"sort"

	"github.com/google/uuid"
	"github.com/mediflow/backend/internal/shared/models"
)

type MatchResult struct {
	DepartmentID uuid.UUID
	Score        float64
	Reason       string
}

// ScorePotentialSource evaluates a potential source department for a request
func ScorePotentialSource(
	request *models.SharingRequest,
	targetDept *models.Department,
	sourceDept *models.Department,
	availableCount int,
	minStock int,
) MatchResult {
	score := 0.0
	reason := ""

	// 1. Proximity Score (Max 50 points)
	if targetDept.Building == sourceDept.Building {
		score += 30
		if targetDept.Floor == sourceDept.Floor {
			score += 20
			reason += "Same building and floor. "
		} else {
			reason += "Same building. "
		}
	} else {
		reason += "Different building. "
	}

	// 2. Stock Health Score (Max 50 points)
	// How many items above min_stock?
	excess := availableCount - minStock
	if excess > 0 {
		score += float64(excess * 10)
		if score > 100 { // Cap at 100 total
			score = 100
		}
		reason += "High stock availability. "
	} else if availableCount > 0 {
		score += 5
		reason += "Limited availability (at minimum). "
	}

	return MatchResult{
		DepartmentID: sourceDept.ID,
		Score:        score,
		Reason:       reason,
	}
}

func FindBestMatches(ctx context.Context, request *models.SharingRequest, candidates []MatchResult) []MatchResult {
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

	if len(candidates) > 3 {
		return candidates[:3]
	}
	return candidates
}
