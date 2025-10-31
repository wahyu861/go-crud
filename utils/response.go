package utils

type BaseResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(message string, data interface{}) BaseResponse {
	return BaseResponse{
		Status:  true,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
}

func ErrorResponse(message string, errs []string) BaseResponse {
	return BaseResponse{
		Status:  false,
		Message: message,
		Errors:  errs,
		Data:    nil,
	}
}