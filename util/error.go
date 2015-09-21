package util

type Error struct {
	Err      interface{}
	Type     string
	IsNormal bool
}

func (err *Error) Error() string {
	return err.Type + ": " + (err.Err.(error)).Error()
}
