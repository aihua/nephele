package request

type DeleteRequest struct {
	IsDeleteAll bool   `xml:"IsDeleteAll"`
	ExPath      string `xml:"ExPath"`
	FilePath    string `xml:"FilePath"`
}
