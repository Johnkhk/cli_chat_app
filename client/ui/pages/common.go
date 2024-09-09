package pages

type errMsg struct {
	err error
}

func (e errMsg) Error() string {
	return e.err.Error()
}
