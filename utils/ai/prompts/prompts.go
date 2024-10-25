package prompts

import (
	"fmt"
)

func FeedbackPrompt(personalStatement string) string {
	return fmt.Sprintf(`Please provide name suggestions for a idea based on the following details:

    You are an expert personal statement reviewer. Please provide detailed feedback on the following personal statement. Your feedback should assess the following aspects:

1. **Clarity**: How clear and concise is the statement? Is it easy to follow?
2. **Structure**: Does the statement have a logical flow? Are ideas presented in a coherent manner?
3. **Grammar & Spelling**: Identify any grammatical errors, spelling mistakes, or awkward phrasing.
4. **Relevance**: Does the content align with the purpose of a personal statement (i.e., showcasing strengths, experiences, and motivation)?
5. **Engagement**: Is the statement compelling? Does it grab attention and hold interest?
6. **Overall Impression**: What is the overall impression of the statement? Is it persuasive and impactful?

Please provide your feedback in the following JSON format:

{
  "clarity": {
    "rating": "<rating out of 10>",
    "feedback": "<feedback text>"
  },
  "structure": {
    "rating": "<rating out of 10>",
    "feedback": "<feedback text>"
  },
  "grammar_spelling": {
    "rating": "<rating out of 10>",
    "feedback": "<feedback text>"
  },
  "relevance": {
    "rating": "<rating out of 10>",
    "feedback": "<feedback text>"
  },
  "engagement": {
    "rating": "<rating out of 10>",
    "feedback": "<feedback text>"
  },
  "overall_impression": {
    "rating": "<rating out of 10>",
    "feedback": "<feedback text>"
  }
}

`, personalStatement)

}
