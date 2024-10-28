package plagiarism

type TextSubmissionResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Text struct {
			ID int `json:"id"`
		} `json:"text"`
	} `json:"data"`
}

type TextStatusResponse struct {
	Data struct {
		State int `json:"state"`
	} `json:"data"`
}

type ReportResponse struct {
	Data struct {
		Report struct {
			Percent string `json:"percent"`
		} `json:"report"`
		ReportData struct {
			Sources []struct {
				Source  string  `json:"source"`
				Percent float64 `json:"percent"`
			} `json:"sources"`
		} `json:"report_data"`
	} `json:"data"`
}

type PlagiarismResult struct {
	PlagiarismPercentage float64
	Sources              []Source
}

type Source struct {
	URL     string
	Percent float64
}
