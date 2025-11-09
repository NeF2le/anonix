package helpers

import (
	errs "github.com/NeF2le/anonix/common/errors"
	"unicode"
)

func ValidatePassword(pw string) error {
	if len(pw) > 72 {
		return errs.ErrPasswordTooLong
	}
	if len(pw) < 8 {
		return errs.ErrPasswordTooShort
	}

	var hasLetter, hasDigit bool

	for _, r := range pw {
		if r > unicode.MaxASCII {
			return errs.ErrPasswordNonASCII
		}

		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit {
		return errs.ErrPasswordWeak
	}

	return nil
}
