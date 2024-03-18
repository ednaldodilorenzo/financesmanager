package util

type ItemNotFoundError struct {
	Message string
}

func (e *ItemNotFoundError) Error() string {
	return e.Message
}
