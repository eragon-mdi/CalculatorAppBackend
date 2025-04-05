package main

import (
	// "encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	// "github.com/google/uuid"
	"net/http"
)

type errorMap map[string]string

type Calculation struct {
	ID         string  `json:"id"`
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
}

type CalculationRequest struct {
	Expression string `json:"expression"`
}

// .
// история выражений
var calculationsHistory = []Calculation{}

// .
// ГЕТ историю
func getCalculations(c echo.Context) error {
	return c.JSON(http.StatusOK, calculationsHistory)
}

// .
// ПОСТ результат вычислений
func postCalculation(c echo.Context) error {
	// получаем запрос????
	var req CalculationRequest
	// декодируем
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalud request"})
	}

	// новый элемент для сохранения в память
	newCalc := Calculation{
		ID:         uuid.NewString(),
		Expression: req.Expression,
	}

	// считаем выражение
	if err := calculatExpression(&newCalc); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}

	// добавляем в историю
	calculationsHistory = append(calculationsHistory, newCalc)
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// ищем старую запись в слайсе
	for i, oldCalc := range calculationsHistory {
		// когда нашли старую запись
		if oldCalc.ID == id {
			// обновляем поле с выражением в старой записи
			var buf string // если нужно будет вернуть старое значение
			buf, calculationsHistory[i].Expression = calculationsHistory[i].Expression, req.Expression

			// заново считаем старую запись
			if err := calculatExpression(&calculationsHistory[i]); err != nil {
				calculationsHistory[i].Expression = buf
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
			}

			// возвращаем успешный ответ
			return c.JSON(http.StatusOK, calculationsHistory[i])
		}
	}
	// если старая запись не была найдена
	return c.JSON(http.StatusBadRequest, errorMap{"error": "Invalid id"})
}

// .
// ДЕЛ выражения по айди
func delCalculation(c echo.Context) error {

	id := c.Param("id")

	var req CalculationRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorMap{"error": "Invalid request"})
	}

	for i, calculation := range calculationsHistory {
		if calculation.ID == id {
			calculationsHistory = append(calculationsHistory[:i], calculationsHistory[i+1:]...)

			return c.NoContent(http.StatusNoContent)
		}
	}
	// если старая запись не была найдена
	return c.JSON(http.StatusBadRequest, errorMap{"error": "Invalid id"})
}

// .
// .
// .
func main() {

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", getCalculations)
	e.POST("/calculations", postCalculation)
	e.PATCH("/calculations/:id", patchCalculation)
	e.DELETE("/calculations/:id", delCalculation)

	e.Start("localhost:8080")

	// json.Marshal()
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
