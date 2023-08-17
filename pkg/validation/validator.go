package validation

type (
	ValidateFunc func(v *Validator)

	Validator struct {
		Errors      []string          `json:",omitempty"`
		FieldErrors map[string]string `json:",omitempty"`
	}
)

func New() *Validator {
	return &Validator{
		Errors:      []string{},
		FieldErrors: map[string]string{},
	}
}

func (v *Validator) HasErrors() bool {
	return len(v.Errors) != 0 || len(v.FieldErrors) != 0
}

func (v *Validator) AddError(message string) {
	if v.Errors == nil {
		v.Errors = []string{}
	}

	v.Errors = append(v.Errors, message)
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = map[string]string{}
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) Check(ok bool, message string) bool {
	if !ok {
		v.AddError(message)
	}

	return ok
}

func (v *Validator) CheckField(ok bool, key, message string) bool {
	if !ok {
		v.AddFieldError(key, message)
	}

	return ok
}

// TODO: add full error description
func (v *Validator) Error() string {
	return "validation error(s)"
}

func Validate(fns ...ValidateFunc) error {
	v := New()

	for _, fn := range fns {
		fn(v)
	}

	if v.HasErrors() {
		return v
	}

	return nil
}
