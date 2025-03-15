package utils

import (
	"github.com/joomcode/errorx"
)

type ErrorCode = string
type Error = *errorx.Error

const (
	errCodeInternal         ErrorCode = "internal"
	errCodeInvalidArgument  ErrorCode = "invalid_argument"
	errCodeNotFound         ErrorCode = "not_found"
	errCodeAlreadyExists    ErrorCode = "already_exists"
	errCodePermissionDenied ErrorCode = "permission_denied"
	errCodeUnauthenticated  ErrorCode = "unauthenticated"
	errCodeUnableToHandle   ErrorCode = "unable_to_handle"
	errCodeInvalidState     ErrorCode = "invalid_state"
	errCodeJSONMarshError   ErrorCode = "json_marsh_error"
	errCodeMultiError       ErrorCode = "multi_error"
)

var (
	ns = errorx.NewNamespace("subswag")

	internal         = errorx.NewType(ns, errCodeInternal)
	invalidArgument  = errorx.NewType(ns, errCodeInvalidArgument)
	notFound         = errorx.NewType(ns, errCodeNotFound)
	alreadyExists    = errorx.NewType(ns, errCodeAlreadyExists)
	permissionDenied = errorx.NewType(ns, errCodePermissionDenied)
	unauthenticated  = errorx.NewType(ns, errCodeUnauthenticated)
	unableToHandle   = errorx.NewType(ns, errCodeUnableToHandle)
	invalidState     = errorx.NewType(ns, errCodeInvalidState)
	jsonMarshError   = errorx.NewType(ns, errCodeJSONMarshError)
	multiError       = errorx.NewType(ns, errCodeMultiError)

	errorTypeMap = map[ErrorCode]errorx.Type{
		errCodeInvalidArgument:  *invalidArgument,
		errCodeNotFound:         *notFound,
		errCodeAlreadyExists:    *alreadyExists,
		errCodePermissionDenied: *permissionDenied,
		errCodeUnauthenticated:  *unauthenticated,
		errCodeUnableToHandle:   *unableToHandle,
		errCodeInvalidState:     *invalidState,
		errCodeJSONMarshError:   *jsonMarshError,
		errCodeMultiError:       *multiError,
	}
)

func NewWrappedError(message string, error error) error {
	if error == nil {
		return nil
	}
	return errorx.DecorateMany(message, error)
}

func NewInternalError(message string, causes ...error) error {
	err := errorx.DecorateMany("cause", causes...)
	return internal.Wrap(err, message)
}

func NewInvalidArgumentError(message string, causes ...error) error {
	err := errorx.DecorateMany("cause", causes...)
	return invalidArgument.Wrap(err, message)
}

func NewNotFoundError(message string, causes ...error) error {
	err := errorx.DecorateMany("cause", causes...)
	return notFound.Wrap(err, message)
}

func NewAlreadyExistsError(message string, causes ...error) error {
	err := errorx.DecorateMany("cause", causes...)
	return alreadyExists.Wrap(err, message)
}

func NewPermissionDeniedError(message string, causes ...error) error {
	err := errorx.DecorateMany("cause", causes...)
	return permissionDenied.Wrap(err, message)
}

func NewUnauthenticatedError(message string, causes ...error) error {
	err := errorx.DecorateMany("cause", causes...)
	return unauthenticated.Wrap(err, message)
}

func NewUnableToHandleError(message string, causes ...error) error {
	err := errorx.DecorateMany("cause", causes...)
	return unableToHandle.Wrap(err, message)
}

func NewInvalidStateError(message string, causes ...error) error {
	err := errorx.DecorateMany("cause", causes...)
	return invalidState.Wrap(err, message)
}

func NewJSONMarshError(message string, causes ...error) error {
	err := errorx.DecorateMany("cause", causes...)
	return jsonMarshError.Wrap(err, message)
}

func NewMultiError(message string, causes ...error) error {
	err := errorx.DecorateMany("cause", causes...)
	return multiError.Wrap(err, message)
}

func GetErrorCode(err error) (ErrorCode, bool) {
	if err == nil {
		return "", false
	}

	if typedErr, ok := err.(*errorx.Error); ok {
		return typedErr.Type().FullName(), true
	}

	return "", false
}

func ErrorOrNil(
	message string,
	fn func(message string, causes ...error) error,
	err error,
) error {
	if err == nil {
		return nil
	}
	return fn(message, err)
}

func IsInternalError(err error) bool {
	return errorx.IsOfType(err, internal)
}

func IsInvalidArgumentError(err error) bool {
	return errorx.IsOfType(err, invalidArgument)
}

func IsNotFoundError(err error) bool {
	return errorx.IsOfType(err, notFound)
}

func IsAlreadyExistsError(err error) bool {
	return errorx.IsOfType(err, alreadyExists)
}

func IsPermissionDeniedError(err error) bool {
	return errorx.IsOfType(err, permissionDenied)
}

func IsUnauthenticatedError(err error) bool {
	return errorx.IsOfType(err, unauthenticated)
}

func IsUnableToHandleError(err error) bool {
	return errorx.IsOfType(err, unableToHandle)
}

func IsInvalidStateError(err error) bool {
	return errorx.IsOfType(err, invalidState)
}

func IsJSONMarshError(err error) bool {
	return errorx.IsOfType(err, jsonMarshError)
}

func IsMultiError(err error) bool {
	return errorx.IsOfType(err, multiError)
}
