package domain

type ValidationError struct {
	Err error
}

func (v *ValidationError) Error() string {
	return v.Err.Error()
}

type NotFoundError struct {
	Err error
}

func (v *NotFoundError) Error() string {
	return v.Err.Error()
}
