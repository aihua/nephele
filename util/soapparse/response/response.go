package response

type Response struct {
	Header         Header         `xml:"Header"`
	SaveResponse    SaveResponse    `xml:"SaveResponse"`
	DeleteResponse  DeleteResponse  `xml:"DeleteResponse"`
	LoadZipResponse LoadZipResponse `xml:"LoadZipResponse"`
	LoadImgResponse LoadImgResponse `xml:"LoadImgResponse"`
}
