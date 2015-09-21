package request

type Header struct {
	UserID          string `xml:"UserID,attr"`
	RequestID       string `xml:"RequestID,attr"`
	RequestType     string `xml:"RequestType,attr"`
	ClientIP        string `xml:"ClientIP,attr"`
	AsyncRequest    bool   `xml:"AsyncRequest,attr"`
	Timeout         int    `xml:"Timeout,attr"`
	MessagePriority int    `xml:"MessagePriority,attr"`
	AssemblyVersion string `xml:"AssemblyVersion,attr"`
	RequestBodySize int    `xml:"RequestBodySize,attr"`
	SerializeMode   string `xml:"SerializeMode,attr"`
	RouteStep       int    `xml:"RouteStep,attr"`
	Culture         string `xml:"Culture,attr"`
	Environment     string `xml:"Environment,attr"`
	ReCallType      string `xml:"ReCallType,attr"`
	ReferenceID     string `xml:"ReferenceID,attr"`
	SessionID       string `xml:"SessionID,attr"`
	SrcAppID        string `xml:"SrcAppID,attr"`
	TransNo         string `xml:"TransNo,attr"`
	UseMemoryQ      bool   `xml:"UseMemoryQ,attr"`
	Version         string `xml:"Version,attr"`
	Via             string `xml:"Via,attr"`
}
