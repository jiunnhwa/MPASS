package errors

type ErrCode string

const ErrCodeUnknown ErrCode = "ERR_UNKNOWN"
const ErrCodeDataNotFound ErrCode = "ERR_DATA_NOT_FOUND"

type Error struct {
	ErrorCode     ErrCode `json:"error_code"`
	InternalError error   `json:"-"`
	Message       string  `json:"message"`
}

func (e Error) Error() string {
	return e.InternalError.Error()
}
