package util

func GetImageConnectionString(tableZone int) string {
	return "gct@tcp(localhost:3306)/test?charset=utf8"
}
