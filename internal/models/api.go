package models

type InfoResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Status string `json:"status"`
	Result Result `json:"result"`
}

type Result struct {
	Output   string `json:"output"`
	ExitCode int    `json:"exit_code"`
}

type ToolManifest struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	InputSchema map[string]any `json:"input_schema,omitempty"`
}

type ManifestResponse struct {
	Tools []ToolManifest `json:"tools"`
}

func NewSuccessResponse(output string, exitCode int) SuccessResponse {
	return SuccessResponse{
		Status: "success",
		Result: Result{
			Output:   output,
			ExitCode: exitCode,
		},
	}
}

func NewErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Error: APIError{
			Code:    code,
			Message: message,
		},
	}
}
