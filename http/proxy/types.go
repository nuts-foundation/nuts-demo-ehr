package proxy

const (
	CodeSecurity  = "security"
	SeverityError = "error"
)

type IssueDetails struct {
	Text string `json:"text"`
}

type Issue struct {
	Code     string        `json:"code"`
	Severity string        `json:"severity"`
	Details  *IssueDetails `json:"details"`
}

type OperationOutcome struct {
	Text  string `json:"text"`
	Issue *Issue `json:"issue"`
}

func NewOperationOutcome(err error, text string, code string, severity string) *OperationOutcome {
	return &OperationOutcome{
		Text: text,
		Issue: &Issue{
			Code:     code,
			Severity: severity,
			Details: &IssueDetails{
				Text: err.Error(),
			},
		},
	}
}
