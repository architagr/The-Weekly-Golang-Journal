package apperrors

type ActiveSessionError struct {
}

func (ActiveSessionError) Error() string {
	return "there is an active session"
}
