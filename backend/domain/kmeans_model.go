package domain

type KMeansModel struct {
	K         int         `json:"k"`
	Features  []string    `json:"features"`
	Centroids [][]float64 `json:"centroids"`
}
