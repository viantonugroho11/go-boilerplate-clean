package apperrors

// Domain sentinel errors. Use errors.Is in repositories and usecases.

var (
	ErrSampleIDRequired  = Validation("sample id is required")
	ErrSampleNotFound    = NotFound("sample not found")
	ErrUserIDRequired    = Validation("id is required")
	ErrUserNotFound      = NotFound("user not found")
	ErrUserNameRequired  = Validation("name is required")
	ErrUserEmailRequired = Validation("email is required")
)
