package request

type Request struct {
	Header         Header         `xml:"Header"`
	SaveRequest    SaveRequest    `xml:"SaveRequest"`
	DeleteRequest  DeleteRequest  `xml:"DeleteRequest"`
	LoadZipRequest LoadZipRequest `xml:"LoadZipRequest"`
	LoadImgRequest LoadImgRequest `xml:"LoadImgRequest"`
}
