package data

type GatorSqlError string

func (e GatorSqlError) Error() string {
	return string(e)
}

const (
	ErrRecordNotFound GatorSqlError = "record not found"
)
