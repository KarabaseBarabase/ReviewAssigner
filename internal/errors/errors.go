package errors

import "fmt"

var (
	ErrPRNotFound     = NewError("NOT_FOUND", "PR not found")
	ErrUserNotFound   = NewError("NOT_FOUND", "User not found")
	ErrTeamNotFound   = NewError("NOT_FOUND", "Team not found")
	ErrPRExists       = NewError("PR_EXISTS", "PR already exists")
	ErrTeamExists     = NewError("TEAM_EXISTS", "Team already exists")
	ErrPRMerged       = NewError("PR_MERGED", "Cannot reassign on merged PR")
	ErrNotAssigned    = NewError("NOT_ASSIGNED", "Reviewer not assigned to this PR")
	ErrNoCandidate    = NewError("NO_CANDIDATE", "No active replacement candidate in team")
	ErrAuthorNotFound = NewError("NOT_FOUND", "Author not found")
)

type Error struct {
	Code    string
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func NewError(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func WrapError(err *Error, cause error) *Error {
	return &Error{
		Code:    err.Code,
		Message: err.Message,
		Cause:   cause,
	}
}

func Is(err error, target *Error) bool {
	if domainErr, ok := err.(*Error); ok {
		return domainErr.Code == target.Code
	}
	return false
}
