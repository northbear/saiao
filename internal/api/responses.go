package api

import "saiao/internal/models"

func SuccessResponse(requestID, output string, exitCode int) models.SuccessResponse {
	return models.SuccessResponse{
		RequestID: requestID,
		Status:    "success",
		Result: models.Result{
			Output:   output,
			ExitCode: exitCode,
		},
	}
}
