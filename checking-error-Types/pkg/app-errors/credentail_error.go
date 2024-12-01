package apperrors

type CredentialError struct {
}

func (CredentialError) Error() string {
	return "invalid credentails"
}
