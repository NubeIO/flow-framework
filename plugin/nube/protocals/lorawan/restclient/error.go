package lwrest

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

// An Error maps to Form3 API error responses
type Error struct {
	Code    int    `json:"error_code,omitempty"`
	Message string `json:"error_message,omitempty"`
}

// Convert error response into error message
func getAPIError(resp *resty.Response) error {
	e := new(Error)
	e.Code = resp.StatusCode()
	e.Message = resp.String()
	return fmt.Errorf("request failed [%d]: %s", e.Code, e.Message)
}
