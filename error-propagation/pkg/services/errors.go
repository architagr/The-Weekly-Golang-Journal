package service

import "error-propagation/pkg/dto"

type ServiceError struct {
	InnerError     error
	HttpStatusCode int
	HttpResponse   *dto.ErrorResponse
}

func (e *ServiceError) Error() string {
	return e.InnerError.Error()
}

func (e *ServiceError) GetHttpStatusCode() int {
	return e.HttpStatusCode
}

func (e *ServiceError) GetErrorResponse() *dto.ErrorResponse {
	return e.HttpResponse
}
