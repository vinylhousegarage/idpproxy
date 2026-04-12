package apierror

type ErrorResponse struct {
	Error ErrorCode `json:"error"`
}
