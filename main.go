package main

import (
	terminHandler "TerminSystem/Handlers/Termin"
	terminService "TerminSystem/Repositories/Termin"
	"TerminSystem/ent"
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/gin-gonic/gin"
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

    r.GET("/termins",TerminHandler.GetAppointmentTimes)

    r.Run(":8080")
}