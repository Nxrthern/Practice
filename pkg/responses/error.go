package responses

type ErrorMessage struct {
	Message string `json:"message"`
	Cause   string `json:"cause"`
}
