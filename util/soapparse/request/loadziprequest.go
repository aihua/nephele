package request

type LoadZipRequest struct {
	Files LoadFileList `xml:"Files"`
}

type LoadFileList struct {
	LoadFiles []LoadFile `xml:"LoadFile"`
}

type LoadFile struct {
	FilePath string `xml:"FilePath"`
	Rename   string `xml:"Rename"`
	IsSource bool   `xml:"IsSource"`
}
