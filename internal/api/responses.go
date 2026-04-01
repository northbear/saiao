package api

import "saiao/internal/models"

func SuccessResponse(output string, exitCode int) models.SuccessResponse {
	return models.SuccessResponse{
		Status: "success",
		Result: models.Result{
			Output:   output,
			ExitCode: exitCode,
		},
	}
}
