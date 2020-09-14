package types

// WebError - This is a common error type used to standardize all api errors
type WebError struct {
	Status     int      `json:"status,omitempty" example:"500"`                                      // status code
	IncidentID string   `json:"incidentId,omitempty" example:"89ed56f9-1935-4e13-9b3c-c208df64b484"` // trace reference id
	Message    string   `json:"message,omitempty" example:"Internal error"`                          // same as http status message
	Errors     []string `json:"errors,omitempty" example:""`                                         // list of actual errors "Field missing required length .. "
	MoreInfo   string   `json:"moreInfo,omitempty" example:""`                                       // link to documentation if any
}
