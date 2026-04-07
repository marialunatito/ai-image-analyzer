package entity

type Tag struct {
	Label      string  `json:"label"`
	Confidence float64 `json:"confidence"`
}

type AnalyzeResult struct {
	Tags []Tag `json:"tags"`
}
