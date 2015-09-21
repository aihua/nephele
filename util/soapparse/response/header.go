package response

type Header struct {
	ServerIP                    string `xml:"ServerIP,attr"`
	ShouldRecordPerformanceTime bool   `xml:"ShouldRecordPerformanceTime,attr"`
	UserID                      string `xml:"UserID,attr"`
	RequestID                   string `xml:"RequestID,attr"`
	ResultCode                  string `xml:"ResultCode,attr"`
	AssemblyVersion             string `xml:"AssemblyVersion,attr"`
	RequestBodySize             int    `xml:"RequestBodySize,attr"`
	SerializeMode               string `xml:"SerializeMode,attr"`
	RouteStep                   int    `xml:"RouteStep,attr"`
	Environment                 string `xml:"Environment,attr"`
	Culture                     string `xml:"Culture,attr"`
	ReferenceID                 string `xml:"ReferenceID,attr"`
	ResultMsg                   string `xml:"ResultMsg,attr"`
	ResultNo                    string `xml:"ResultNo,attr"`
	SessionID                   string `xml:"SessionID,attr"`
	Timestamp                   string `xml:"Timestamp,attro"`
	TransNo                     string `xml:"TransNo,attr"`
	Version                     string `xml:"Version.attr"`
}
