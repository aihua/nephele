package util

type Error struct {
	Err      interface{}
	Type     string
	IsNormal bool
}
