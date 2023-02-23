package twitter

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

var (
	UsernameMinLength = 2
	PasswordMinLength = 6
	emailRegexp       = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (AuthResponse, error)
	Login(ctx context.Context, input LoginInput) (AuthResponse, error)
}

type AuthToken struct {
	ID  string
	Sub string
}

type AuthResponse struct {
	AccessToken string
	User        User
}

type RegisterInput struct {
	Email           string
	Username        string
	Password        string
	ConfirmPassword string
}

func (ri *RegisterInput) Sanitize() {
	ri.Email = strings.TrimSpace(ri.Email)
	ri.Email = strings.ToLower(ri.Email)

	ri.Username = strings.TrimSpace(ri.Username)
}

func (ri RegisterInput) Validate() error {
	if len(ri.Username) < UsernameMinLength {
		return fmt.Errorf("%w: username not long enough, (%d) characters at least", ErrValidation, UsernameMinLength)
	}

	if !emailRegexp.MatchString(ri.Email) {
		return fmt.Errorf("%w: email not valid", ErrValidation)
	}

	if len(ri.Password) < PasswordMinLength {
		return fmt.Errorf("%w: password not long enough, (%d) characters at least", ErrValidation, PasswordMinLength)
	}

	if ri.Password != ri.ConfirmPassword {
		return fmt.Errorf("%w: confirm password must match the password", ErrValidation)
	}

	return nil
}

type LoginInput struct {
	Email    string
	Password string
}

func (li *LoginInput) Sanitize() {
	li.Email = strings.TrimSpace(li.Email)
	li.Email = strings.ToLower(li.Email)
}

func (li LoginInput) Validate() error {
	if !emailRegexp.MatchString(li.Email) {
		return fmt.Errorf("%w: email not valid", ErrValidation)
	}

	if len(li.Password) < 1 {
		return fmt.Errorf("%w: password required", ErrValidation)
	}

	return nil
}
