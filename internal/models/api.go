package models

type InfoResponse struct {
	RequestID string `json:"request_id,omitempty"`
	Status    string `json:"status"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
}

type ErrorResponse struct {
	RequestID string   `json:"request_id,omitempty"`
	Error     APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	RequestID string `json:"request_id,omitempty"`
	Status    string `json:"status"`
	Result    Result `json:"result"`
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
	RequestID string         `json:"request_id,omitempty"`
	Format    string         `json:"format,omitempty"`
	Tools     []ToolManifest `json:"tools"`
}

type OpenAITool struct {
	Type        string         `json:"type"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters"`
	Strict      bool           `json:"strict"`
}

type OpenAIManifestResponse struct {
	RequestID string       `json:"request_id,omitempty"`
	Format    string       `json:"format"`
	Tools     []OpenAITool `json:"tools"`
}

func NewSuccessResponse(requestID, output string, exitCode int) SuccessResponse {
	return SuccessResponse{
		RequestID: requestID,
		Status:    "success",
		Result: Result{
			Output:   output,
			ExitCode: exitCode,
		},
	}
}

func NewErrorResponse(requestID, code, message string) ErrorResponse {
	return ErrorResponse{
		RequestID: requestID,
		Error: APIError{
			Code:    code,
			Message: message,
		},
	}
}
