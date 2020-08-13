package exception

type Exception struct {
	errors []error
}

func New() *Exception {
	return &Exception{}
}

func (e *Exception) Put(err error) {
	e.errors = append(e.errors, err)
}

func (e *Exception) First() error {
	if len(e.errors) > 0 {
		return e.errors[0]
	}
	return nil
}

func (e *Exception) Errors() []error {
	return e.errors
}

func (e *Exception) Has() bool {
	if len(e.errors) > 0 {
		return true
	}
	return false
}

func (e *Exception) Reset() {
	if len(e.errors) > 0 {
		e.errors = []error{}
	}
}

func (e *Exception) Clone() *Exception {
	newObj := &Exception{}
	if len(e.errors) > 0 {
		for _, err := range e.errors {
			newObj.Put(err)
		}
	}
	return newObj
}
