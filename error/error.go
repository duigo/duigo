package error

type Error struct {
	Id     int
	Desc   string
	Locale string
	Params map[string]string
	Resaon *Error
}
