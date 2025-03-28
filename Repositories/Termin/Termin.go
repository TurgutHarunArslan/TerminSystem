package termin

import (
	"TerminSystem/ent"
	"TerminSystem/ent/appointment"
	"context"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"time"
)

type AppointmentService struct {
	client *ent.Client
}

func NewAppointmentService(client *ent.Client) *AppointmentService {
	return &AppointmentService{
		client: client,
	}
}

func (s *AppointmentService) GetBusinessHours(weekday time.Weekday) (int, int) {
	switch weekday {
	case time.Saturday:
		return 10, 14
	case time.Sunday:
		return -1, -1
	default:
		return 10, 17
	}
}

func (s *AppointmentService) IsValidTerminDate(dateStr string, timeStr string, now ...time.Time) (bool, error) {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false, InvalidDateError(dateStr)
	}

	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return false, LocationLoadError()
	}

	var currentTime time.Time
	if len(now) > 0 && !now[0].IsZero() {
		currentTime = now[0]
	} else {
		currentTime = time.Now().In(loc)
	}

	var targetTime time.Time
	if timeStr != "" {
		parsedTime, err := time.Parse("15:04", timeStr)
		if err != nil {
			return false, InvalidDateError(timeStr)
		}

		loc, err := time.LoadLocation("Europe/Berlin")
		if err != nil {
			return false, LocationLoadError()
		}

		targetTime = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, loc)
	} else {
		targetTime = parsedDate
	}

	weekday := parsedDate.Weekday()
	startHour, endHour := s.GetBusinessHours(weekday)
	if startHour == -1 {
		return false, DateShopClosedError(weekday.String())
	}

	dateOnly := timeStr == ""
	var compareTime time.Time

	if dateOnly {
		compareTime = targetTime.Truncate(24 * time.Hour)
		currentDate := currentTime.Truncate(24 * time.Hour)
		if compareTime.Before(currentDate) {
			return false, DateInPastError(compareTime.String(), currentDate.String())
		}
	} else {
		if targetTime.Before(currentTime) {
			return false, DateInPastError(targetTime.String(), currentTime.String())
		}
	}

	if timeStr != "" && (targetTime.Hour() < startHour || targetTime.Hour() >= endHour) {
		return false, DateShopClosedError(targetTime.String())
	}

	return true, nil
}

func (s *AppointmentService) GetAvailableDates(ctx context.Context, days int) []string {
	var availableDates []string

	loc, err := time.LoadLocation("Europe/Berlin")

	if err != nil {
		fmt.Println(err.Error())
		return availableDates
	}

	today := time.Now().In(loc).Truncate(24 * time.Hour)
	createdDays := 0
	for i := 0; createdDays < days; i++ {
		date := today.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")

		isValid, _ := s.IsValidTerminDate(dateStr, "")
		if !isValid {
			continue
		}

		availableDates = append(availableDates, dateStr)
		createdDays += 1
	}
	return availableDates
}

func (s *AppointmentService) GetTimeSlotsByDate(ctx context.Context, dateStr string) ([]string, error) {
	parsedDate, err := time.Parse("2006-01-02", dateStr)

	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return nil, LocationLoadError()
	}

	currentTime := time.Now().In(loc)

	isValid, err := s.IsValidTerminDate(dateStr, "")

	if !isValid {
		return nil, err
	}

	if err != nil {
		return nil, InvalidDateError(dateStr)
	}

	var terminSlots []string
	if parsedDate.Before(currentTime.Truncate(24 * time.Hour)) {
		return nil, DateInPastError(dateStr, currentTime.String())
	}

	weekday := currentTime.Weekday()

	startHour, endHour := s.GetBusinessHours(weekday)

	for hour := startHour; hour < endHour; hour++ {
		for minute := 0; minute < 60; minute += 30 {
			timeStr := fmt.Sprintf("%02d:%02d", hour, minute)
			isValid, err := s.IsValidTerminDate(parsedDate.Format("2006-01-02"), timeStr, currentTime)
			if !isValid || err != nil {
				continue
			}

			datestr := parsedDate.Add(time.Hour*time.Duration(hour) + time.Minute*time.Duration(minute))
			terminSlots = append(terminSlots, datestr.Format("2006-01-02 15:04"))
		}
	}

	return terminSlots, nil
}

func (s *AppointmentService) BookAppointment(ctx context.Context, name, email, phone, desc string, Type appointment.Type, date time.Time) (*ent.Appointment, error) {
	delkey, err := gonanoid.New(128)
	if err != nil {
		return nil, err
	}


	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return nil, LocationLoadError()
	}

	if time.Now().In(loc).Add(24 * time.Hour * 28).Before(date) {
		return nil, DateNotReadyError(date.String())
	}

	isValid, err := s.IsValidTerminDate(date.Truncate(24*time.Hour).Format("2006-01-02"), date.Format("15:04"))
	if !isValid || err != nil {
		return nil, err
	}

	return s.client.Appointment.Create().
		SetName(name).
		SetEmail(email).
		SetPhone(phone).
		SetStartTime(date).
		SetEndTime(date.Add(30 * time.Minute)).
		SetDescription(desc).
		SetType(Type).
		SetDelkey(delkey).
		Save(ctx)
}

func (s *AppointmentService) DeleteAppointment(ctx context.Context, delkey string) error {
	_, err := s.client.Appointment.Delete().Where(appointment.DelkeyEQ(delkey)).Exec(ctx)
	return err
}