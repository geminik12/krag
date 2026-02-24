package validation

import (
	"context"
	"errors"
	"regexp"

	v1 "github.com/geminik12/krag/pkg/api/apiserver/v1"
)

// ValidateCreateUserRequest 校验 CreateUserRequest
func (v *Validator) ValidateCreateUserRequest(ctx context.Context, rq *v1.CreateUserRequest) error {
	// Validate username
	if rq.Username == "" {
		return errors.New("username cannot be empty")
	}
	if len(rq.Username) < 4 || len(rq.Username) > 32 {
		return errors.New("username must be between 4 and 32 characters")
	}

	// Validate password
	if rq.Password == "" {
		return errors.New("password cannot be empty")
	}
	if len(rq.Password) < 8 || len(rq.Password) > 64 {
		return errors.New("password must be between 8 and 64 characters")
	}

	// Validate nickname (if provided)
	if rq.Nickname != nil && *rq.Nickname != "" {
		if len(*rq.Nickname) > 32 {
			return errors.New("nickname cannot exceed 32 characters")
		}
	}

	// Validate email
	if rq.Email == "" {
		return errors.New("email cannot be empty")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`).MatchString(rq.Email) {
		return errors.New("email format is invalid")
	}

	// Validate phone number
	if rq.Phone == "" {
		return errors.New("phone number cannot be empty")
	}
	if !regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(rq.Phone) {
		return errors.New("phone number format is invalid")
	}

	return nil
}

// ValidateUpdateUserRequest 校验 UpdateUserRequest
func (v *Validator) ValidateUpdateUserRequest(ctx context.Context, rq *v1.UpdateUserRequest) error {
	if rq.Username != nil {
		if len(*rq.Username) < 4 || len(*rq.Username) > 32 {
			return errors.New("username must be between 4 and 32 characters")
		}
	}

	if rq.Nickname != nil {
		if len(*rq.Nickname) > 32 {
			return errors.New("nickname cannot exceed 32 characters")
		}
	}

	if rq.Email != nil {
		if !regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`).MatchString(*rq.Email) {
			return errors.New("email format is invalid")
		}
	}

	if rq.Phone != nil {
		if !regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(*rq.Phone) {
			return errors.New("phone number format is invalid")
		}
	}

	return nil
}

// ValidateLoginRequest 校验 LoginRequest
func (v *Validator) ValidateLoginRequest(ctx context.Context, rq *v1.LoginRequest) error {
	if rq.Username == "" {
		return errors.New("username cannot be empty")
	}
	if rq.Password == "" {
		return errors.New("password cannot be empty")
	}
	return nil
}

// ValidateChangePasswordRequest 校验 ChangePasswordRequest
func (v *Validator) ValidateChangePasswordRequest(ctx context.Context, rq *v1.ChangePasswordRequest) error {
	if rq.OldPassword == "" {
		return errors.New("old password cannot be empty")
	}
	if rq.NewPassword == "" {
		return errors.New("new password cannot be empty")
	}
	if len(rq.NewPassword) < 8 || len(rq.NewPassword) > 64 {
		return errors.New("new password must be between 8 and 64 characters")
	}
	if rq.OldPassword == rq.NewPassword {
		return errors.New("new password cannot be the same as old password")
	}
	return nil
}

// ValidateListUserRequest 校验 ListUserRequest
func (v *Validator) ValidateListUserRequest(ctx context.Context, rq *v1.ListUserRequest) error {
	if rq.Limit < 0 {
		return errors.New("limit must be greater than or equal to 0")
	}
	if rq.Offset < 0 {
		return errors.New("offset must be greater than or equal to 0")
	}
	if rq.Limit > 100 {
		return errors.New("limit cannot exceed 100")
	}
	return nil
}
