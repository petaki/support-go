package forms

// Error type.
type Error struct {
	Status int
	Msg    string
}

// Error function.
func (e *Error) Error() string {
	return e.Msg
}
