package schemes

import "strconv"

// ValidateError --------------------------------
type ValidateError struct {
	Field string
	Msg   string
}

type ValidateErrorResponse struct {
	Code  int             `json:"code"`
	Error []ValidateError `json:"error"`
}

type ErrorResponse struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e *ErrorResponse) Error() string {
	return "ERROR: " + e.Err + "(Code " + strconv.Itoa(e.Code) + ")"
}
