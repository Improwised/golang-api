package utils

import "regexp"

func ValidateEmail(email string) (bool, error) {
	return regexp.MatchString("[a-zA-z]+@improwised.com", email)
}
