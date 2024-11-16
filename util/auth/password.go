package auth

import (
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

const (
	PasswordMinLength               = 12
	PasswordSpecialChars            = "`-~!@#$%^&*()_+='\",./?><"
	PasswordRequireNumbers          = true
	PasswordRequireSpecialChars     = true
	PasswordRequireCapitalLetters   = true
	PasswordRequireLowercaseLetters = true
)

const (
	pwdRemoveSpacesPattern    = `\s+`
	pwdNumberPattern          = `\d+`
	pwdSpecialCharsPattern    = `[` + PasswordSpecialChars + `]+`
	pwdCapitalLetterPattern   = `[A-Z]+`
	pwdLowercaseLetterPattern = `[a-z]+`

	pwdHashCostMultiple = 6
)

type PasswordOptions struct {
	MinLength               int
	RequireSpecialChars     bool
	RequireNumbers          bool
	RequireCapitalLetters   bool
	RequireLowercaseLetters bool
}

var (
	passwordOptions = PasswordOptions{
		MinLength:               PasswordMinLength,
		RequireSpecialChars:     PasswordRequireSpecialChars,
		RequireNumbers:          PasswordRequireNumbers,
		RequireCapitalLetters:   PasswordRequireCapitalLetters,
		RequireLowercaseLetters: PasswordRequireLowercaseLetters,
	}

	sRe  = regexp.MustCompile(pwdRemoveSpacesPattern)
	nRe  = regexp.MustCompile(pwdNumberPattern)
	scRe = regexp.MustCompile(pwdSpecialCharsPattern)
	clRe = regexp.MustCompile(pwdCapitalLetterPattern)
	llRe = regexp.MustCompile(pwdLowercaseLetterPattern)
)

func PasswordMeetsRequirements(password string) bool {
	if sRe.MatchString(password) {
		return false
	}

	if len(password) < passwordOptions.MinLength {
		return false
	}

	if passwordOptions.RequireNumbers {
		if !nRe.MatchString(password) {
			return false
		}
	}

	if passwordOptions.RequireSpecialChars {
		if !scRe.MatchString(password) {
			return false
		}
	}

	if passwordOptions.RequireCapitalLetters {
		if !clRe.MatchString(password) {
			return false
		}
	}

	if passwordOptions.RequireLowercaseLetters {
		if !llRe.MatchString(password) {
			return false
		}
	}

	return true
}

func HashPassword(password string) (string, error) {
	hash, hErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost*pwdHashCostMultiple)

	if hErr != nil {
		return "", hErr
	}

	return string(hash), nil
}

func PasswordMatches(password string, pwdHash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(password)) == nil
}

func SetPasswordOptions(options PasswordOptions) {
	passwordOptions = options
}

func GetPasswordOptions() PasswordOptions {
	return passwordOptions
}
