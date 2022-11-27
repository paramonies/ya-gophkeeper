package core

import (
	"errors"
	"fmt"
)

// UniqueViolationError is err, that should be used when data couldn't be inserted
// due to unsatisfied condition of unique value.
type UniqueViolationError struct {
	Constraint string
	Message    string
}

// NewUniqueViolationError returns new UniqueViolationError. It accepts name of unique constraint.
func NewUniqueViolationError(constraint, message string) error {
	return &UniqueViolationError{
		Constraint: constraint,
		Message:    message,
	}
}

func (e *UniqueViolationError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return fmt.Sprintf("violation of unique constraint %s", e.Constraint)
}

// IsUniqueViolationError returns true if the err's chain contains UniqueViolationError error.
// Otherwise, it returns false.
func IsUniqueViolationError(err error) bool {
	var want *UniqueViolationError
	return errors.As(err, &want)
}

type ResourceNotFoundError struct {
	Kind string
	ID   string
}

func (e *ResourceNotFoundError) Error() string {
	return fmt.Sprintf("%s %q not found", e.Kind, e.ID)
}

type UserNotFoundError struct {
	ResourceNotFoundError
}

func NewUserNotFoundError(login string) error {
	return &UserNotFoundError{
		ResourceNotFoundError{
			Kind: "user",
			ID:   login,
		},
	}
}

func (e *UserNotFoundError) Unwrap() error {
	return &e.ResourceNotFoundError
}

type PasswordNotFoundError struct {
	ResourceNotFoundError
}

func NewPasswordNotFoundError(login string) error {
	return &PasswordNotFoundError{
		ResourceNotFoundError{
			Kind: "password",
			ID:   login,
		},
	}
}

func (e *PasswordNotFoundError) Unwrap() error {
	return &e.ResourceNotFoundError
}

type TextNotFoundError struct {
	ResourceNotFoundError
}

func NewTextNotFoundError(title string) error {
	return &TextNotFoundError{
		ResourceNotFoundError{
			Kind: "text",
			ID:   title,
		},
	}
}

func (e *TextNotFoundError) Unwrap() error {
	return &e.ResourceNotFoundError
}

type BinaryNotFoundError struct {
	ResourceNotFoundError
}

func NewBinaryNotFoundError(title string) error {
	return &BinaryNotFoundError{
		ResourceNotFoundError{
			Kind: "binary",
			ID:   title,
		},
	}
}

func (e *BinaryNotFoundError) Unwrap() error {
	return &e.ResourceNotFoundError
}

type CardNotFoundError struct {
	ResourceNotFoundError
}

func NewCardNotFoundError(number string) error {
	return &CardNotFoundError{
		ResourceNotFoundError{
			Kind: "card",
			ID:   number,
		},
	}
}

func (e *CardNotFoundError) Unwrap() error {
	return &e.ResourceNotFoundError
}

// IsNotFound returns true if the err's chain contains ResourceNotFound error.
// Otherwise, it returns false.
func IsNotFound(err error) bool {
	var want *ResourceNotFoundError
	return errors.As(err, &want)
}
