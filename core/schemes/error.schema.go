package schemes

import "strconv"

// ValidateError --------------------------------
type ValidateError struct {
	Field string
	Msg   string
}

type ShmValidateErrorResponse struct {
	Code  int             `json:"code"`
	Error []ValidateError `json:"error"`
}

type ShmErrorResponse struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e *ShmErrorResponse) Error() string {
	return "ERROR: " + e.Err + "(Code " + strconv.Itoa(e.Code) + ")"
}
