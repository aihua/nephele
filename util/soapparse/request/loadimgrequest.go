package request

type LoadImgRequest struct {
	IsSource bool   `xml:"IsSource"`
	FilePath string `xml:"FilePath"`
}
