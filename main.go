package main

import (
	terminHandler "TerminSystem/Handlers/Termin"
	terminService "TerminSystem/Repositories/Termin"
	"TerminSystem/ent"
	"TerminSystem/templates"
	"context"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)


func main() {
    ctx := context.Background()
    client, err := ent.Open("sqlite3", "file:appointment.db?mode=rwc&_fk=1")
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer client.Close()

    if err := client.Schema.Create(ctx); err != nil {
        log.Fatalf("Failed to create schema: %v", err)
    }

    TerminService := terminService.NewAppointmentService(client)
    TerminHandler := terminHandler.NewTerminHandle(TerminService)

    r := gin.Default()
    gin.SetMode(gin.DebugMode)

    api := r.Group("/api")

    api.GET("/termins",TerminHandler.GetAppointmentTimes)
    api.POST("/termins",TerminHandler.BookAppoinment)
    api.DELETE("/termins",TerminHandler.DeleteAppoinment)

    r.GET("/", func(c *gin.Context) {
        c.Status(http.StatusOK)
        c.Header("Content-Type","text/html")
        templates.Root().Render(c.Request.Context(),c.Writer)
    }) 

    r.Run(":8080")
}