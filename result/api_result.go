package result

type RestfulResult[T any] struct {
	Data    T      `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
	Success bool   `json:"success"`
}

func NewSuccessResult[T any](data T, message string) RestfulResult[T] {
	return RestfulResult[T]{
		Data:    data,
		Message: message,
		Code:    200,
		Success: true,
	}
}

func NewSuccessResultWithCode[T any](data T, code int, message string) RestfulResult[T] {
	if code < 200 || code >= 300 {
		panic("Success code must be in the range 200-299")
	}
	return RestfulResult[T]{
		Data:    data,
		Message: message,
		Code:    code,
		Success: true,
	}
}

func NewErrorResult[T any](message string, code int) RestfulResult[T] {
	return RestfulResult[T]{
		Message: message,
		Code:    code,
		Success: false,
	}
}
