package presenters

import "github.com/gofiber/fiber/v2"

type BatchOperation struct {
	OperationType        string `json:"operationType"`
	ItemType             string `json:"itemType"`
	NumberOfItemAffected int64  `json:"numberOfItemAffected"`
}

func ErrorResponse(msg string) *fiber.Map {
	return &fiber.Map{
		"status":  "error",
		"message": msg,
	}
}

func FailResponse(data map[string]interface{}) *fiber.Map {
	return &fiber.Map{
		"status": "fail",
		"data":   data,
	}
}

func BatchOperationResponse(operationType, itemType string, numberOfItemAffected int64) *fiber.Map {
	batchOps := &BatchOperation{
		OperationType:        operationType,
		ItemType:             itemType,
		NumberOfItemAffected: numberOfItemAffected,
	}
	return &fiber.Map{
		"status": "success",
		"data":   batchOps,
	}
}
