package validation

type Validator struct {
	Errors      []string          `json:"errors,omitempty"`
	FieldErrors map[string]string `json:"fieldErrors,omitempty"`
}

func New() *Validator {
	return &Validator{}
}

func (v *Validator) HasErrors() bool {
	return len(v.Errors) != 0 || len(v.FieldErrors) != 0
}

func (v *Validator) AddError(message string) *Validator {
	if v.Errors == nil {
		v.Errors = []string{}
	}

	v.Errors = append(v.Errors, message)

	return v
}

func (v *Validator) AddFieldError(key, message string) *Validator {
	if v.FieldErrors == nil {
		v.FieldErrors = map[string]string{}
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}

	return v
}

func (v *Validator) Check(ok bool, message string) {
	if !ok {
		v.AddError(message)
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func (v *Validator) Error() string {
	// TODO: Build a better error message
	return "validation error(s)"
}

func (v *Validator) As(target any) bool {
	_, ok := target.(*Validator)
	return ok
}

func Validate(fn func(v *Validator)) error {
	v := New()
	fn(v)
	if v.HasErrors() {
		return v
	}
	return nil
}
