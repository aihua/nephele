package response

type SaveResponse struct {
	CheckPass    bool    `xml:"CheckPass"`
	OriginalPath string  `xml:"OriginalPath"`
	TargetPath   string  `xml:"TargetPath"`
	Process      Process `xml:"Process"`
}

type Process struct {
	ProcessResponses []ProcessResponse `xml:"ProcessResponse"`
}

type ProcessResponse struct {
	ID   string `xml:"ID"`
	Path string `xml:"Path"`
	Info string `xml:"Info"`
}
