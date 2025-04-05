package math

import (
	"errors"
	"fmt"
	"strconv"
)

// в такой структуре храним выражение в виде среза операндов и среза операций (операторов)
type expression struct {
	operands  []float64
	operators []rune
}

// Минимально нужная структура для "не занания"
type Calculable interface {
	GetExpression() string
	SetResult(float64)
}

// принимает указатель на структуру
// считает результат, иначе возвращает ошибку
// ТОЛЬКО бинарные операторы!
// func CalculateExpression(c *Calculation) (err error) {
func CalculateExpression(c Calculable) error {
	// парсим
	exp, err := parsing(c.GetExpression())
	if err != nil {
		return err
	}

	// первое значение
	result := exp.operands[0]

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
		result = operationFunc(result, exp.operands[i+1])

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

	c.SetResult(result)
	return nil
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
