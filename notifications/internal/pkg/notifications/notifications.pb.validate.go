// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: notifications.proto

package notifications

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on GetHistoryRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetHistoryRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetHistoryRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetHistoryRequestMultiError, or nil if none found.
func (m *GetHistoryRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetHistoryRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if val := m.GetUserId(); val < 1 || val > 9223372036854775807 {
		err := GetHistoryRequestValidationError{
			field:  "UserId",
			reason: "value must be inside range [1, 9223372036854775807]",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetStartTime()) != 22 {
		err := GetHistoryRequestValidationError{
			field:  "StartTime",
			reason: "value length must be 22 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)

	}

	if !_GetHistoryRequest_StartTime_Pattern.MatchString(m.GetStartTime()) {
		err := GetHistoryRequestValidationError{
			field:  "StartTime",
			reason: "value does not match regex pattern \"^\\\\d{4}-\\\\d{2}-\\\\d{2} \\\\d{2}:\\\\d{2}:\\\\d{2}\\\\+\\\\d{2}$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetEndTime()) != 22 {
		err := GetHistoryRequestValidationError{
			field:  "EndTime",
			reason: "value length must be 22 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)

	}

	if !_GetHistoryRequest_EndTime_Pattern.MatchString(m.GetEndTime()) {
		err := GetHistoryRequestValidationError{
			field:  "EndTime",
			reason: "value does not match regex pattern \"^\\\\d{4}-\\\\d{2}-\\\\d{2} \\\\d{2}:\\\\d{2}:\\\\d{2}\\\\+\\\\d{2}$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GetHistoryRequestMultiError(errors)
	}

	return nil
}

// GetHistoryRequestMultiError is an error wrapping multiple validation errors
// returned by GetHistoryRequest.ValidateAll() if the designated constraints
// aren't met.
type GetHistoryRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetHistoryRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetHistoryRequestMultiError) AllErrors() []error { return m }

// GetHistoryRequestValidationError is the validation error returned by
// GetHistoryRequest.Validate if the designated constraints aren't met.
type GetHistoryRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetHistoryRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetHistoryRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetHistoryRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetHistoryRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetHistoryRequestValidationError) ErrorName() string {
	return "GetHistoryRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetHistoryRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetHistoryRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetHistoryRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetHistoryRequestValidationError{}

var _GetHistoryRequest_StartTime_Pattern = regexp.MustCompile("^\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}\\+\\d{2}$")

var _GetHistoryRequest_EndTime_Pattern = regexp.MustCompile("^\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}\\+\\d{2}$")

// Validate checks the field values on GetHistoryResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetHistoryResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetHistoryResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetHistoryResponseMultiError, or nil if none found.
func (m *GetHistoryResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetHistoryResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetEntries() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetHistoryResponseValidationError{
						field:  fmt.Sprintf("Entries[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetHistoryResponseValidationError{
						field:  fmt.Sprintf("Entries[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetHistoryResponseValidationError{
					field:  fmt.Sprintf("Entries[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return GetHistoryResponseMultiError(errors)
	}

	return nil
}

// GetHistoryResponseMultiError is an error wrapping multiple validation errors
// returned by GetHistoryResponse.ValidateAll() if the designated constraints
// aren't met.
type GetHistoryResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetHistoryResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetHistoryResponseMultiError) AllErrors() []error { return m }

// GetHistoryResponseValidationError is the validation error returned by
// GetHistoryResponse.Validate if the designated constraints aren't met.
type GetHistoryResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetHistoryResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetHistoryResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetHistoryResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetHistoryResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetHistoryResponseValidationError) ErrorName() string {
	return "GetHistoryResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetHistoryResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetHistoryResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetHistoryResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetHistoryResponseValidationError{}

// Validate checks the field values on Entry with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Entry) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Entry with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in EntryMultiError, or nil if none found.
func (m *Entry) ValidateAll() error {
	return m.validate(true)
}

func (m *Entry) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for OrderId

	// no validation rules for Status

	// no validation rules for CreatedAt

	if len(errors) > 0 {
		return EntryMultiError(errors)
	}

	return nil
}

// EntryMultiError is an error wrapping multiple validation errors returned by
// Entry.ValidateAll() if the designated constraints aren't met.
type EntryMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m EntryMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m EntryMultiError) AllErrors() []error { return m }

// EntryValidationError is the validation error returned by Entry.Validate if
// the designated constraints aren't met.
type EntryValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e EntryValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e EntryValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e EntryValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e EntryValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e EntryValidationError) ErrorName() string { return "EntryValidationError" }

// Error satisfies the builtin error interface
func (e EntryValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sEntry.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = EntryValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = EntryValidationError{}
