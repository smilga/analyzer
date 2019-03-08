package api

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrScaning = errors.New("Error scanning crypt string value")
)

// Cryptstring bcrypts given string when saved to sql
type Cryptstring string

// Scan to implement Scanner interface for sql
func (c *Cryptstring) Scan(value interface{}) error {
	if bs, err := driver.String.ConvertValue(value); err == nil {
		if str, ok := bs.([]byte); ok {
			*c = Cryptstring(str)
			return nil
		}
		return fmt.Errorf("%s: %v", ErrScaning, value)
	}
	return fmt.Errorf("%s: %v", ErrScaning, value)
}

// Value to implemen valuer interface for sql
func (c Cryptstring) Value() (driver.Value, error) {
	ok, err := IsBcrypt(string(c))
	if err != nil {
		return nil, err
	}
	if ok {
		return string(c), nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(c), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return string(hash), nil
}

//IsBcrypt checks whether string is bcrypt hash
func IsBcrypt(s string) (bool, error) {
	reg, err := regexp.Compile(`^\$2[ayb]\$.{56}$`)
	if err != nil {
		return false, err
	}
	ok := reg.MatchString(s)
	if !ok {
		return false, nil
	}
	return true, nil
}
