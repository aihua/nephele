package models

type ImagePlan struct {
	ImgIdx       int64
	Plan         string
	PartitionKey int16
	TableZone    int
}
