package disgolink

var _ error = (*Exception)(nil)

func NewExceptionFromErr(err error) *Exception {
	return &Exception{Message: err.Error(), Severity: SeverityFault}
}

func NewException(message string, severity Severity) *Exception {
	return &Exception{Message: message, Severity: severity}
}

type Exception struct {
	Message  string   `json:"message"`
	Severity Severity `json:"severity"`
	Cause    *string  `json:"cause,omitempty"`
}

func (e *Exception) Error() string {
	return e.Message
}
