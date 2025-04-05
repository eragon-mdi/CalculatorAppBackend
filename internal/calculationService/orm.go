package calculationservice

import "kalc/internal/calculationService/math"

type Calculation struct {
	ID         string  `json:"id" gorm:"primaryKey"`
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
}

type CalculationRequest struct {
	Expression string `json:"expression"`
}

// .
func (c *Calculation) CalculatExpression() (err error) {
	return math.CalculateExpression(c)
}

func (c *Calculation) GetExpression() string {
	return c.Expression
}
func (c *Calculation) SetResult(res float64) {
	c.Result = res
}
