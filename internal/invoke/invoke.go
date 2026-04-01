package invoke

import "saiao/internal/models"

func Execute(_ string, _ string, _ []byte, _ *models.Config) (models.SuccessResponse, error) {
	return models.SuccessResponse{}, nil
}
