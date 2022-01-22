package lavalink

var _ error = (*FriendlyException)(nil)

type FriendlyException struct {
	Message  string   `json:"message"`
	Severity Severity `json:"severity"`
	Cause    *string  `json:"cause,omitempty"`
}

func (e FriendlyException) Error() string {
	return e.Message
}

type Severity string

//goland:noinspection GoUnusedConst
const (
	SeverityCommon     Severity = "COMMON"
	SeveritySuspicious Severity = "SUSPICIOUS"
	SeverityFault      Severity = "FAULT"
)
