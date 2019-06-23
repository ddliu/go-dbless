package dbless

type RecordNotFoundError struct {
}

func (e RecordNotFoundError) Error() string {
	return "Record not found"
}

func IsRecordNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(RecordNotFoundError)

	return ok
}
