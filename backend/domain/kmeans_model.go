package domain

type KMeansModel struct {
	K               int         `json:"k"`
	Features        []string    `json:"features"`
	Centroids       [][]float64 `json:"centroids"`
	Means           []float64   `json:"means"`
	Stds            []float64   `json:"stds"`
	AvgTargetDeltas []float64   `json:"avg_target_deltas"`
}

type KMeansFeatures struct {
	TargetDelta      float64
	HasBrokerage     float64
	ActionScore      float64
	RatingDeltaScore float64
	TimeDelta        float64
}

func (f KMeansFeatures) FeaturesToSlice() []float64 {
	return []float64{
		f.TargetDelta,
		f.HasBrokerage,
		f.ActionScore,
		f.RatingDeltaScore,
		f.TimeDelta,
	}
}
