package files

type MandatoryError struct {
	field string
}

func (e *MandatoryError) Error() string {
	return e.field + " is mandatory"
}

type InvalidInputError struct {
}

func (e *InvalidInputError) Error() string {
	return "Invalid input"
}
