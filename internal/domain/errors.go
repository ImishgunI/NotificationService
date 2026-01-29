package domain

type BusinessError struct {
	Reason string
}

type RetryableError struct {
	Err error
}

type InfrasractureError struct {
	Err error
}

func (e BusinessError) Error() string {
	return e.Reason
}

func (r RetryableError) Error() string {
	return r.Err.Error()
} 

func (r RetryableError) Unwrap() error {
	return r.Err
}

func (i InfrasractureError) Error() string {
	return i.Err.Error()
}

func (i InfrasractureError) Unwrap() error {
	return i.Err
}

