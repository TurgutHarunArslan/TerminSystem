package termin

import (
	"fmt"
)

const (
	InvalidDateErrorCode = iota
	DateInPastErrorCode
	DateShopClosedErrorCode
	DateNotReadyErrorCode
	LocationLoadErrorCode
)

type AppointmentError struct {
	Code    int
	Message string
	Details string
}

func (e *AppointmentError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Details: %s", e.Code, e.Message, e.Details)
}

func NewAppointmentError(code int, message, details string) *AppointmentError {
	return &AppointmentError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// InvalidDateError creates an error for invalid date format with dynamic message
func InvalidDateError(dateStr string) error {
	return NewAppointmentError(InvalidDateErrorCode, "invalid date format", "Invalid date string: "+dateStr)
}

// DateInPastError creates an error when the date is in the past with dynamic message
func DateInPastError(pastDate, currentDate string) error {
	return NewAppointmentError(DateInPastErrorCode, "appointment date is in the past", "Target date: "+pastDate+" is before current date " +currentDate)
}

func DateShopClosedError(dateStr string) error {
	return NewAppointmentError(DateShopClosedErrorCode, "appointment date is out of working time", "Target date: "+dateStr+" is after working time")
}

func DateNotReadyError(dateStr string) error {
	return NewAppointmentError(DateNotReadyErrorCode, "appointment date is out of working time", "Target date: "+dateStr+" is after working time")
}

func LocationLoadError() error {
	return NewAppointmentError(LocationLoadErrorCode,"Location Load Failed","Could not load  Berlin time")
}