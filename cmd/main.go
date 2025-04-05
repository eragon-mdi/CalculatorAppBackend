package main

import (
	"log"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	calculationservice "kalc/internal/calculationService"
	db "kalc/internal/db"
	"kalc/internal/handlers"
)

func main() {

	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	// слои
	calcRepo := calculationservice.NewCalculationRepository(database)
	calcService := calculationservice.NewCalcService(calcRepo)
	calcHandlers := handlers.NewCalculationHandler(calcService)

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", calcHandlers.GetCalculations)
	e.POST("/calculations", calcHandlers.PostCalculation)
	e.PATCH("/calculations/:id", calcHandlers.PatchCalculation)
	e.DELETE("/calculations/:id", calcHandlers.DelCalculation)

	e.Start("localhost:8080")
}
