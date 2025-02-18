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
	Code    int
	Err     string
	ErrBase error
}

func (e *HTTPError) Error() string {
	return "#" + e.Err + "(Code " + strconv.Itoa(e.Code) + ")\n"
}

type HTTPError struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e *ErrorResponse) Error() string {
	if e.ErrBase == nil {
		return "ErrorResponse == nil"
	}
	return "" + e.Err + "(Code " + strconv.Itoa(e.Code) + ")  // " + e.ErrBase.Error()
}
