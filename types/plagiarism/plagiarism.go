package plagiarism

type PlagiarismResult struct {
	Score   float64  `json:"score"`
	Sources []Source `json:"sources"`
}

type Source struct {
	URL   string  `json:"url"`
	Score float64 `json:"score"`
}
