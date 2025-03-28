package termin

import (
	termin "TerminSystem/Repositories/Termin"
	"TerminSystem/ent/appointment"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

type TerminHandler struct {
	service *termin.AppointmentService
}

func NewTerminHandle(service *termin.AppointmentService) *TerminHandler {
	return &TerminHandler{
		service: service,
	}
}

func (h *TerminHandler) GetAppointmentTimes(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datum ist erforderlich"})
		return
	}

	times, err := h.service.GetTimeSlotsByDate(c.Request.Context(), date)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if len(times) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Keine Termine"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": times})
	return
}

type AppoinmentCreate struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Desc  string `json:"desc"`
	Type  string `json:"type"`
	Date  string `json:"date"`
}

func (h *TerminHandler) BookAppoinment(c *gin.Context) {
	var CreateData AppoinmentCreate
	if err := c.ShouldBind(&CreateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date,err := time.Parse("2006-01-02 15:04",CreateData.Date)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appoinment, err := h.service.BookAppointment(c.Request.Context(),CreateData.Name,CreateData.Email,CreateData.Phone,CreateData.Desc,appointment.Type(CreateData.Type),date)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK,gin.H{"data":appoinment.String()})
	return
}



func (h *TerminHandler) DeleteAppoinment(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key ist erforderlich"})
		return
	}

	err := h.service.DeleteAppointment(c.Request.Context(), key)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Termin Losen"})
	return
}