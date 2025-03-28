package termin

import (
	"TerminSystem/ent/appointment"
	"TerminSystem/ent/enttest"
	"context"
	"strings"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func isWeekday(date time.Time) bool {
	weekday := date.Weekday()
	return weekday != time.Saturday && weekday != time.Sunday
}

func TestBookAppointment(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		t.Fatalf("failed creating schema: %v", err)
	}

	ctx := context.Background()

	service := NewAppointmentService(client)

	dates := service.GetAvailableDates(ctx, 14)
	assert.NotEmpty(t, dates)

	var day time.Time
	for i, v := range dates {
		if i == 0 {
			continue
		}
		date, err := time.Parse("2006-01-02", v)
		assert.NoError(t, err)

		if isWeekday(date) {
			day = date
			break
		}
	}

	name := "Test User"
	email := "example@example.com"
	phone := "123456789"
	description := "Test Appointment"
	Type := appointment.TypeSonstiges
	start := day.Add(10 * time.Hour)

	appointment, err := service.BookAppointment(ctx, name, email, phone, description, Type, start)
	assert.NoError(t, err)
	assert.NotNil(t, appointment)

	assert.Len(t, strings.Split(appointment.Delkey, ""), 128)

	appointments, err := client.Appointment.Query().All(ctx)
	assert.NoError(t, err)
	assert.Len(t, appointments, 1)
	assert.Equal(t, name, appointments[0].Name)

	assert.WithinDuration(t, start, appointments[0].StartTime, time.Second)
}

func TestGetAvailableDatesLength(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		t.Fatalf("failed creating schema: %v", err)
	}

	ctx := context.Background()
	service := NewAppointmentService(client)

	dates := service.GetAvailableDates(ctx, 14)
	assert.Len(t, dates, 14)
}

func TestWeekday(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		t.Fatalf("failed creating schema: %v", err)
	}

	ctx := context.Background()
	service := NewAppointmentService(client)

	dates := service.GetAvailableDates(ctx, 7)

	var tomorrow string
	for i, v := range dates {
		if i == 0 {
			continue
		}
		date, err := time.Parse("2006-01-02", v)
		assert.NoError(t, err)

		if isWeekday(date) {
			tomorrow = v
			break
		}
	}

	timeslots, err := service.GetTimeSlotsByDate(ctx, tomorrow)

	assert.NoError(t, err)
	assert.NotNil(t, timeslots)

	TimeStart := timeslots[0]
	TimeEnd := timeslots[len(timeslots)-1]

	assert.Equal(t, tomorrow+" 10:00", TimeStart)
	assert.Equal(t, tomorrow+" 16:30", TimeEnd)
}

func TestDateValidator(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		t.Fatalf("failed creating schema: %v", err)
	}

	service := NewAppointmentService(client)

	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatal("Could not load Loc")
	}

	today := time.Now().In(loc).Truncate(24 * time.Hour)

	tests := []struct {
		name      string
		date      string
		time      string
		reference time.Time
		expected  bool
		expectErr bool
		ErrId     int
	}{
		{
			name:      "past date",
			date:      today.AddDate(0, 0, -1).Format("2006-01-02"),
			time:      "",
			reference: today,
			expected:  false,
			expectErr: true,
			ErrId:     DateInPastErrorCode,
		},
		{
			name:      "today",
			date:      today.Format("2006-01-02"),
			time:      "",
			reference: today,
			expected:  true,
			expectErr: false,
			ErrId:     -1,
		},
		{
			name:      "valid date and time",
			date:      today.Format("2006-01-02"),
			time:      today.Add(11 * time.Hour).Format("15:04"),
			reference: today,
			expected:  true,
			expectErr: false,
			ErrId:     -1,
		},
		{
			name:      "valid date without time (edgecase to check if reference date gets compared without time)",
			date:      today.Format("2006-01-02"),
			time:      "",
			reference: today,
			expected:  true,
			expectErr: false,
			ErrId:     -1,
		},
		{
			name:      "next Saturday",
			date:      today.AddDate(0, 0, (6-int(today.Weekday())+7)%7).Format("2006-01-02"),
			time:      "",
			reference: today,
			expected:  true,
			expectErr: false,
			ErrId:     -1,
		},
		{
			name:      "invalid time on Saturday",
			date:      today.AddDate(0, 0, (6-int(today.Weekday())+7)%7).Format("2006-01-02"),
			time:      today.Add(15 * time.Hour).Format("15:04"),
			reference: today,
			expected:  false,
			expectErr: true,
			ErrId:     DateShopClosedErrorCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid, err := service.IsValidTerminDate(tt.date, tt.time, tt.reference)

			assert.Equal(t, tt.expected, isValid)
			if tt.expectErr {
				assert.Error(t, err)

				customErr, ok := err.(*AppointmentError)
				if !ok {
					t.Fatal("Wrong Error type return?", err.Error())
				}
				if tt.ErrId > -1 {
					assert.Equal(t, tt.ErrId, customErr.Code)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
