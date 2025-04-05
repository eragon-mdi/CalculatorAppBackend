package main

import (
	// "encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"net/http"

	"github.com/google/uuid"

	// .
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type tinyErrJson map[string]string

type Calculation struct {
	ID         string  `json:"id" gorm:"primaryKey"`
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
}

type CalculationRequest struct {
	Expression string `json:"expression"`
}

// подкючаемся к БД
func initDB() {
	dsn := "host=localhost user=postgres password=pass dbname=postgres port=5432 sslmode=disable" // data source name
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) // пустой конфиг

	if err != nil {
		log.Fatal("Не смог подключиться к БД: ", err)
	}

	// автомиграция
	if err := db.AutoMigrate(&Calculation{}); err != nil {
		log.Fatal("Поломался на миграции: ", err)
	}
}

// .
// ГЕТ историю
func getCalculations(c echo.Context) error {
	var calculationsHistory []Calculation

	if res := db.Find(&calculationsHistory); res.Error != nil {
		c.JSON(http.StatusInternalServerError, tinyErrJson{"error": "Invalid db request"})
	}

	fmt.Println("\n", calculationsHistory, "\n")

	return c.JSON(http.StatusOK, calculationsHistory)
}

// .
// ПОСТ результат вычислений
func postCalculation(c echo.Context) error {
	// получаем запрос????
	var req CalculationRequest
	// декодируем
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tinyErrJson{"error": "Invalud request"})
	}

	// новый элемент для сохранения в память
	newCalc := Calculation{
		ID:         uuid.NewString(),
		Expression: req.Expression,
	}

	// считаем выражение
	if err := calculatExpression(&newCalc); err != nil {
		return c.JSON(http.StatusBadRequest, tinyErrJson{"error": "Invalid expression"})
	}

	// добавляем в историю
	// calculationsHistory = append(calculationsHistory, newCalc)
	if res := db.Create(&newCalc); res.Error != nil {
		return c.JSON(http.StatusInternalServerError, tinyErrJson{"error": "Couldn't create new calculation in db"})
	}

	// возвращаем ответ
	return c.JSON(http.StatusCreated, newCalc)
}

// .
// ПАТЧ выражения по айди
func patchCalculation(c echo.Context) error {

	// айдишник из заголовка (ссылки)
	id := c.Param("id")

	// .
	var req CalculationRequest
	// .
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tinyErrJson{"error": "Invalid request"})
	}

	// новая запись
	var newCalc = Calculation{
		Expression: req.Expression,
	}
	// считаем результат нового выражения
	if err := calculatExpression(&newCalc); err != nil {
		return c.JSON(http.StatusBadRequest, tinyErrJson{"error": "Invalid expression"})
	}

	// обновляем запись по айди
	res := db.Model(&Calculation{}).Where("id = ?", id).Updates(newCalc)
	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, tinyErrJson{"error": "Couldn't update the record"})
	}

	// .
	return c.JSON(http.StatusOK, newCalc)
}

// .
// ДЕЛ выражения по айди
func delCalculation(c echo.Context) error {

	id := c.Param("id")

	var req CalculationRequest
	// биндим
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tinyErrJson{"error": "Invalid request"})
	}

	// удаялем по айди
	if res := db.Delete(&Calculation{}, "id = ?", id); res.Error != nil {
		return c.JSON(http.StatusInternalServerError, tinyErrJson{"error": "Couldnt delete from bd or bad id"})
	}

	// .
	return c.NoContent(http.StatusNoContent)
}

// .
// .
// .
func main() {

	initDB()

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", getCalculations)
	e.POST("/calculations", postCalculation)
	e.PATCH("/calculations/:id", patchCalculation)
	e.DELETE("/calculations/:id", delCalculation)

	e.Start("localhost:8080")
}

// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// Далее математика

// в такой структуре храним выражение в виде среза операндов и среза операндов
type expression struct {
	operands  []float64
	operators []rune
}

// принимает указатель на структуру
// считает результат, иначе возвращает ошибку
// ТОЛЬКО бинарные операторы!
func calculatExpression(c *Calculation) (err error) {
	// парсим
	exp, err := parsing(c.Expression)
	if err != nil {
		return err
	}

	// первое значение
	c.Result = exp.operands[0]

	fmt.Printf("Log:\t%+v\n", exp)

	// мапа функций по руне (выёбываюсь)
	operations := map[rune]func(a, b float64) float64{
		'+': func(a, b float64) float64 { return a + b },
		'-': func(a, b float64) float64 { return a - b },
		'*': func(a, b float64) float64 { return a * b },
		'/': func(a, b float64) float64 { return a / b },
	}

	// перебираем и считаем
	for i, operator := range exp.operators { // идём по срезу операторов коих на 1 меньше чем операндов
		operationFunc, exists := operations[operator]
		if !exists {
			return errors.New("Не известный оператор")
		}
		c.Result = operationFunc(c.Result, exp.operands[i+1])

		//	switch operator {
		// 	case '+':
		// 		c.Result += exp.operands[i+1]
		// 	case '-':
		// 		c.Result -= exp.operands[i+1]
		// 	case '*':
		// 		c.Result *= exp.operands[i+1]
		// 	case '/':
		// 		c.Result /= exp.operands[i+1]
		// 	default:
		// 		return errors.New("Не известный оператор")
		//	}
	}

	return
}

// .
func parsing(input string) (res expression, err error) {

	// сначала операнд
	if !('0' <= input[0] && input[0] <= '9') {
		return res, errors.New("Выражение начаинается не с операнда")
	}
	// в конце также должен быть операнд
	if !('0' <= input[len(input)-1] && input[len(input)-1] <= '9') {
		return res, errors.New("Выражение заканчивается не операндом")
	}

	var buf = "" // буферная для res.operands[i]
	// var buf float64

	// парсим
	for _, char := range input {
		switch char {
		case ' ':
			continue
		// операнд
		case '.':
			if buf[len(buf)-1] == '.' { // две точки к ряду
				return res, errors.New("Что за двоеточие ..")
			}
			fallthrough
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			buf += string(char)
		// операции
		case '+', '-', '*', '/':
			// не может быть двух операций подряд, между ними должен быть "не пустой" операнд
			if buf == "" || buf == "." {
				return res, errors.New("Два оператора к ряду")
			}

			// преобразуем строку буфер в float64
			newOperand, err := strconv.ParseFloat(buf, 64)
			if err != nil {
				return res, err
			}
			buf = "" //

			// обновляем стркутуру res
			res.operands = append(res.operands, newOperand) // добавляем число которое до этого собирали
			res.operators = append(res.operators, char)     // также добавляем операцию
		default:
			return res, errors.New("Не известный оператор")
		}
	}

	// последний операнд
	newOperand, err := strconv.ParseFloat(buf, 64)
	if err != nil {
		return res, err
	}
	res.operands = append(res.operands, newOperand)

	return
}
