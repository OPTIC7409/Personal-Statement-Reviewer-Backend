package plagiarismdetection

import (
	"psr/types/plagiarism"
)

func CheckPlagiarism(text string) (*plagiarism.PlagiarismResult, error) {
	// Since we don't have access to the API yet, return fake data
	return &plagiarism.PlagiarismResult{
		PlagiarismPercentage: 23.5,
		Sources: []plagiarism.Source{
			{
				URL:     "https://example.com/source1",
				Percent: 12.5,
			},
			{
				URL:     "https://example.com/source2",
				Percent: 11.0,
			},
		},
	}, nil
}
