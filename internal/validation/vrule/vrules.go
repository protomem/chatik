package vrule

import "github.com/protomem/chatik/internal/validation"

func Nickname(nickname string) validation.ValidateFunc {
	return func(v *validation.Validator) {
		v.CheckField(validation.Length(nickname, 4, 18), "nickname", "must be between 4 and 20 characters")
	}
}

func Password(password string) validation.ValidateFunc {
	return func(v *validation.Validator) {
		v.CheckField(validation.Length(password, 6, 20), "password", "must be between 6 and 18 characters")
	}
}

func Email(email string) validation.ValidateFunc {
	return func(v *validation.Validator) {
		v.CheckField(validation.Email(email), "email", "must be a valid email address")
	}
}
