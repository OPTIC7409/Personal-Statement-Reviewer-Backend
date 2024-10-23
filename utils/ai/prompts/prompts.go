package prompts

import (
	"fmt"
)

func NameGen(ideaName, ideaSummary, idealCustomer, monitisationSummary, nameFormatInstruction string) string {
	prompt := fmt.Sprintf(`Please provide name suggestions for a idea based on the following details:

    Idea Name: '%s'
    Idea Summary: '%s'
    Ideal Customer: '%s'
    Monetization Summary: '%s'

Provide 10 creative and relevant name suggestions for the idea, taking into account the summary, ideal customer, and monetization plan. Ensure that the names are catchy, memorable, and suitable for the idea's market and audience.

%s

The response must be in the following format:
{
"names": ["", "", "", "", "", "", "", "", "", ""]
}

Ensure the JSON object is on a new line by itself.
`, ideaName, ideaSummary, idealCustomer, monitisationSummary, nameFormatInstruction)

	return prompt
}

func SummariseTrendData(trendDataJson []byte) string {
	prompt := fmt.Sprintf(`
	Analyze the following trend data and provide a summary based on the scores for each date. The summary should include assessments for current, future, and past demand.
	
	Here is the trend data:
	%s
	

	Please format the summary as follows:
	{
	   "in_demand_currently": "<true/false>",
	   "in_demand_currently_desc": "<description of current demand>",
	   "in_future_demand": "<true/false>",
	   "in_future_demand_desc": "<description of future demand>",
	   "past_demand": "<true/false>",
	   "past_demand_desc": "<description of past demand>"
	}
	
	DO MAKE SURE THAT ALL DATA IS IN STRING FORMAT.
	YOU MUST MAKE SURE YOU RETURN JSON IN THE FORMAT ABOVE.
	YOU MUST MAKE SURE THAT IT IS VALID JSON CODE BEFORE RETURNING IT.
	YOU SHOULD NOT RETURN ME CODE OR ANYTHING OTHER THAN THE JSON OBJECT.
	
	Make sure to base your analysis on the trend scores, indicating whether the keyword is in demand currently, in the future, and in the past, along with appropriate descriptions for each aspect.
	`, trendDataJson)
	return prompt
}

func GenerateSWOTAnalysis(textData string) string {
	prompt := fmt.Sprintf(`
Analyze the following data and generate a SWOT analysis with at least 3 points for each category: Strengths, Weaknesses, Opportunities, and Threats. Provide detailed explanations for each point.

Data: %s

Format the response as JSON:
{
	"strengths": ["Strength 1", "Strength 2", "Strength 3"],
	"weaknesses": ["Weakness 1", "Weakness 2", "Weakness 3"],
	"opportunities": ["Opportunity 1", "Opportunity 2", "Opportunity 3"],
	"threats": ["Threat 1", "Threat 2", "Threat 3"]
}

Only return the JSON object. No additional text.
`, textData)

	fmt.Println(prompt)
	return prompt
}

func GenerateSWOTAnalysisInDepth(textData string) string {
	prompt := fmt.Sprintf(`
Conduct a comprehensive SWOT analysis on the provided data. Your analysis should be thorough, insightful, and cover at least three points for each category: Strengths, Weaknesses, Opportunities, and Threats. Each point should be supported by detailed explanations and logical reasoning based on the data provided.

Data: %s

Instructions:
- **Strengths**: Identify the core advantages and positive attributes inherent in the data. Consider aspects such as unique features, competitive advantages, and internal resources that contribute to success.
- **Weaknesses**: Highlight the internal limitations or areas of improvement. These could include resource constraints, gaps in capabilities, or any factors that may hinder performance.
- **Opportunities**: Explore external factors or trends that could be leveraged for growth and success. Consider market trends, technological advancements, or changes in consumer behavior that align with the data.
- **Threats**: Identify potential external challenges or risks that could impact the data negatively. Consider competitive pressures, regulatory changes, or economic shifts that may pose a threat.

Format the response as a JSON object:
{
	"strengths": ["Strength 1", "Strength 2", "Strength 3"],
	"weaknesses": ["Weakness 1", "Weakness 2", "Weakness 3"],
	"opportunities": ["Opportunity 1", "Opportunity 2", "Opportunity 3"],
	"threats": ["Threat 1", "Threat 2", "Threat 3"]
}

Ensure that the JSON object is valid and contains no additional text. The analysis should be concise yet comprehensive, providing a clear understanding of each category based on the data.
`, textData)

	return prompt
}

func GenerateFinancialProjections(projectName, projectSummary, idealCustomer, monetisationSummary string) string {
	prompt := fmt.Sprintf(`
    Please provide a detailed financial projection for the following project based on the given details:
    
    Project Name: '%s'
    Project Summary: '%s'
    Ideal Customer: '%s'
    Monetization Summary: '%s'
    
    Include the following details:
    - Initial investment cost breakdown (e.g., server costs, logo design, marketing, etc.)
    - Monthly operating expenses (e.g., hosting, salaries, utilities, etc.)
    - Revenue forecasting for the first year, broken down by month
    - Potential profitability calculation
    - A detailed explanation of the financial projection

    The detailed explanation should include the assumptions made, the methodology used, and any other relevant information.

    It should be noted that the financial projection should be as detailed as possible and should be based on the information provided above.

    Also note that the financial projection should be realistic and achievable based on the information provided.

    PLEASE ENSURE ALL FINANCIAL FIGURES ARE IN USD.
    PLEASE ENSURE ALL DATA YOU PROVIDE IS IN STRINGS. FOR EXAMPLE: ("$xxx,xxx")

    DO NOT PROVIDE ANY OTHER INFORMATION APART FROM THE FINANCIAL PROJECTION.

    DO NOT PUT ANY ELLIPSES ("...") IN THE RESPONSE. PROVIDE FULL LISTINGS FOR EACH SECTION. 

    PLEASE MAKE SURE ANY NAMES WITH SPACES USE AN UNDERSCORE INSTEAD OF A SPACE. FOR EXAMPLE: "OPERATING EXPENSES" SHOULD BE "OPERATING_EXPENSES".
    ALSO PLEASE MAKE SURE THAT YOU DO NOT JUST PROVIDE A SINGLE FIGURE AND INCREMENT IT FOR THE REVENUE FORECASTING.
    YOU MUST ALSO BE REALISTIC WITH THE REVENUE FORECASTING AND HONEST AS TO THE POTENTIAL PROFITABILITY OF THE PROJECT.

    PLEASE BARE IN MIND THAT THESE ARE START UP PROJECTS AND SHOULD BE TREATED AS SUCH. THEY SHOULD NOT HAVE $100,000 INITIAL INVESTMENTS OR $1,000,000 REVENUE FORECASTS. AS THIS IS INCREADIBLY UNREALISTIC. 

    ENSURE THE REVENUE FORECASTING INCLUDES DATA FOR EACH OF THE 12 MONTHS OF THE FIRST YEAR.

    YOU CAN PROVIDE AS MANY COSTS AS YOU NEED, BUT YOU MUST ENSURE THAT YOU ARE PROVIDING A DETAILED FINANCIAL PROJECTION.

The response must be in the following format:
{
    "monthly": {
        "cost_name": "$x",
        "another_cost_name": "$x"
    },
    "yearly": {
        "cost_name": "$x",
        "another_cost_name": "$x"
    },
    "one_time": {
        "cost_name": "$x",
        "another_cost_name": "$x"
    }
}

Ensure the JSON object is on a new line by itself.
`, projectName, projectSummary, idealCustomer, monetisationSummary)

	return prompt
}

func GenerateCommentsForTrends(responseDataJson []byte) string {
	prompt := fmt.Sprintf(`
    Analyze the following trend data and provide one positive and one negative comment.
    
    Trend Data:
    %s
    
    Please format the response as follows:
    {
        "positive_comment": "<positive comment>",
        "negative_comment": "<negative comment>"
    }
    
    Ensure that you dont mention words that hint it being a "dataset" or "data".
    You can mention dates and times, but not in a way that hints it being a dataset, where you are just saying "on this date there was a spike in demand" etc instead of saying 'over time'.
    Ensure the comments are insightful and based on the trend data provided.
    Each comment should be around 10 words long.
    `, responseDataJson)
	return prompt
}

func PredictFutureOfTrend(trendDataJson []byte, duration int) string {
	fmt.Println("PredictFutureOfTrend: Trend data:", string(trendDataJson))
	prompt := fmt.Sprintf(`
    Analyze the following trend data and provide insights on both positive and negative aspects of the data. Additionally, predict the future of the trend with a minimum of 20 separate points based on the previous data.

    Trend Data:
    %s

    The data you are given is a list of dates and their corresponding trend scores for the idea that we are analysing.

    Please format the response as follows:
    {
        "comments": {
            "positive_comment": "<positive comment>",
            "negative_comment": "<negative comment>"
        },
        "future_trend": [
            {"point": 1, "value": "<predicted value>"},
            {"point": 2, "value": "<predicted value>"},
            ...
            {"point": 30, "value": "<predicted value>"}
        ]
    }

    The duration of the trend data is over %d days.
    You should provide a prediction for this amount of time while still returning 30 points.
        
    You can mention dates and times, but not in a way that hints it being a dataset, where you are just saying "on this date there was a spike in demand" etc instead of saying 'over time'.
    Ensure the comments are insightful and the future trend predictions are based on the trend data provided.
    Each comment should be around 10 words long.
    `, trendDataJson, duration)
	return prompt
}
