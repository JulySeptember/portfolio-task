package httpx

// =========================
// error codes
// =========================

const (

	// common
	CodeInternalServerError = "INTERNAL_SERVER_ERROR"
	CodeValidationError     = "VALIDATION_ERROR"
	CodeInvalidJSON         = "INVALID_JSON"
	CodeInvalidID           = "INVALID_ID"
	CodeMethodNotAllowed    = "METHOD_NOT_ALLOWED"

	// auth
	CodeUnauthorized = "UNAUTHORIZED"
	CodeInvalidToken = "INVALID_TOKEN"

	// user
	CodeUserNotFound   = "USER_NOT_FOUND"
	CodeDuplicateEmail = "DUPLICATE_EMAIL"

	// task
	CodeTaskNotFound   = "TASK_NOT_FOUND"
	CodeInvalidUserID  = "INVALID_USER_ID"
	CodeInvalidDueDate = "INVALID_DUE_DATE"
	CodeInvalidSort    = "INVALID_SORT"
	CodeInvalidOrder   = "INVALID_ORDER"
)

// =========================
// common error response
// =========================

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// =========================
// validation error response
// =========================

type ValidationErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}
