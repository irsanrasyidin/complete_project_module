package exception

type Exception struct {
	Code     Code
	Message  string
	err      error
	details  any
	httpCode int
}

type Option func(*Exception)

func New(code Code, message string, opts ...Option) *Exception {
	exc := &Exception{
		Code:     code,
		Message:  message,
		httpCode: HTTPStatus(code),
	}

	for _, opt := range opts {
		opt(exc)
	}

	return exc
}

func WithError(err error) Option {
	return func(exc *Exception) {
		exc.err = err
	}
}

func WithDetails(details any) Option {
	return func(exc *Exception) {
		exc.details = details
	}
}

func WithHTTPCode(httpCode int) Option {
	return func(exc *Exception) {
		exc.httpCode = httpCode
	}
}

func (exc *Exception) Error() string {
	if exc == nil {
		return ""
	}
	if exc.err != nil {
		return exc.Message + ": " + exc.err.Error()
	}
	return exc.Message
}

func (exc *Exception) Unwrap() error {
	if exc == nil {
		return nil
	}
	return exc.err
}

func (exc *Exception) GetError() any {
	if exc == nil {
		return nil
	}
	if exc.details != nil {
		return exc.details
	}
	if exc.err != nil {
		return exc.err.Error()
	}
	return nil
}

func (exc *Exception) GetHttpCode() int {
	if exc == nil {
		return HTTPStatus(InternalErrorCode)
	}
	if exc.httpCode > 0 {
		return exc.httpCode
	}
	return HTTPStatus(exc.Code)
}
