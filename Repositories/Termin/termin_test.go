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

	dates := service.GetAvailableDates(ctx,14)
	assert.NotEmpty(t,dates)


	var day time.Time
	for _, v := range dates {
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

	appointment, err := service.BookAppointment(ctx,name, email, phone, description, Type, start)
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

	var tomorrow  string
	for _, v := range dates {
		date, err := time.Parse("2006-01-02", v)
		assert.NoError(t, err)

		if isWeekday(date) {
			tomorrow = v
			break
		}
	}

	timeslots, err := service.GetTimeSlotsByDate(ctx,tomorrow)

	assert.NoError(t, err)
	assert.NotNil(t, timeslots)

	TimeStart := timeslots[0]
	TimeEnd := timeslots[len(timeslots)-1]

	assert.Equal(t, tomorrow + " 10:00", TimeStart)
	assert.Equal(t, tomorrow + " 16:30", TimeEnd)
}


func TestDateValidator(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		t.Fatalf("failed creating schema: %v", err)
	}

	service := NewAppointmentService(client)


	today := time.Now().Truncate(24 *time.Hour)

	isValid,err := service.IsValidTerminDate(today.AddDate(0,0,-1).Format("2006-01-02"),"")

	assert.False(t,isValid)
	assert.Error(t,err)

	isValid,err = service.IsValidTerminDate(today.Format("2006-01-02"),today.Add(11 *time.Hour).Format("15:04"),today)

	assert.True(t,isValid)
	assert.NoError(t,err)

	daysUntilSaturday := (6 - int(today.Weekday()) + 7) % 7
	if daysUntilSaturday == 0 {
		daysUntilSaturday = 7
	}

	isValid,err = service.IsValidTerminDate(time.Now().AddDate(0,0,daysUntilSaturday).Format("2006-01-02"),"",today)

	assert.True(t,isValid)
	assert.NoError(t,err)


	isValid,err = service.IsValidTerminDate(time.Now().AddDate(0,0,daysUntilSaturday).Format("2006-01-02"),today.Add(15 * time.Hour).Format("15:04"),today)

	assert.False(t,isValid)
	assert.Error(t,err)
}