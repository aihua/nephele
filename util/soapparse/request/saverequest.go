package request

type SaveRequest struct {
	FileBytes      string      `xml:"FileBytes"`
	FilePath       string      `xml:"FilePath"`
	CheckItem      CheckItem   `xml:"CheckItem"`
	OriginalFormat string      `xml:"OriginalFormat"`
	SaveType       string      `xml:"SaveType"`
	TargetFormat   string      `xml:"TargetFormat"`
	SourceSaveMode string      `xml:"SourceSaveMode"`
	TargetQuality  int         `xml:"TargetQuality"`
	ObjectID       string      `xml:"ObjectID"`
	Channel        string      `xml:"Channel"`
	StorageType    string      `xml:"StorageType"`
	Process        ProcessList `xml:"Process"`
	SourcePath     string      `xml:"SourcePath"`
	PlanID         string      `xml:"PlanID"`
	StartDateTime  string      `xml:"StartDateTime"`
	PathType       int         `xml:"PathType"`
}

type CheckItem struct {
	IsOtherImage bool       `xml:"IsOtherImage"`
	Types        FormatList `xml:"Types"`
	MaxBytes     float64    `xml:"MaxBytes"`
	MinWidth     int        `xml:"MinWidth"`
	MinHeight    int        `xml:"MinHeight"`
	MaxWidth     int        `xml:"MaxWidth"`
	MaxHeight    int        `xml:"MaxHeight"`
}

type FormatList struct {
	ImgFormats []string `xml:"ImgFormat"`
}

type AnyType struct {
	//ProcessGroup
	Type          string      `xml:"type,attr"`
	ID            int         `xml:"ID"`
	ExName        string      `xml:"ExName"`
	TargetFormat  string      `xml:"TargetFormat"`
	TargetQuality string      `xml:"TargetQuality"`
	IsSharpen     bool        `xml:"IsSharpen"`
	Radius        int         `xml:"Radius"`
	Sigma         int         `xml:"Sigma"`
	Channels      string      `xml:"Channels"`
	List          ProcessList `xml:"List"`
	//Cropping
	Proportion float64 `xml:"Proportion"`
	//Watermark
	Transparency float64 `xml:"Transparency"`
	Position     string  `xml:"Position"`
	X            int     `xml:"X"`
	Y            int     `xml:"Y"`
	MinWidth     int     `xml:"MinWith"` //As it's wrong in C#, we have to be wrong
	MinHeight    int     `xml:"MinHeight"`
	FileBytes    []byte  `xml:"FileBytes"`
	//Compression
	Width  int `xml:"Width"`
	Height int `xml:"Height"`
	//Conver
	//ProportionCompression
	//CustomCropping
	//SpecialGroup
	Taipe int `xml:"Type"`
}

type ProcessList struct {
	AnyTypes []AnyType `xml:"anyType"`
}
