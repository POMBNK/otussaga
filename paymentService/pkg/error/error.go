package error

import (
	"encoding/json"
	"fmt"
)

var ErrNotFound = New("Not found", "", "US-000404")

type ApiError struct {
	Description  string `json:"description"`
	Err          error  `json:"-"`
	DeveloperMsg string `json:"developer_msg"`
	Code         string `json:"code"`
}

func New(description, developerMsg, code string) *ApiError {
	return &ApiError{
		Description:  description,
		Err:          fmt.Errorf(description),
		DeveloperMsg: developerMsg,
		Code:         code,
	}
}

func (e *ApiError) Unwrap() error { return e.Err }

func (e *ApiError) Error() string {
	return e.Err.Error()
}

func (e *ApiError) Marshal() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil
	}

	return bytes
}

func systemErr(developerMsg string) *ApiError {
	return New("system error", developerMsg, "US-000500")
}
