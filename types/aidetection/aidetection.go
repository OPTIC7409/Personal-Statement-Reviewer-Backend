package aidetection

type AIDetectionResult struct {
	OverallAIProbability float64          `json:"overall_ai_probability"`
	FlaggedSections      []FlaggedSection `json:"flagged_sections"`
}

type FlaggedSection struct {
	Text        string  `json:"text"`
	Reason      string  `json:"reason"`
	Probability float64 `json:"probability"`
}
