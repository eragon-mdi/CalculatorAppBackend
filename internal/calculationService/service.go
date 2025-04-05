package calculationservice

import (
	"github.com/google/uuid"
)

// CGUD
type CalculationService interface {
	Creaters
	Geters
	Updaters
	Deleters
}

type Geters interface {
	AllCalculationsGeter
	CalculationGeter
}

type Creaters interface {
	CreateCalculation(exp string) (Calculation, error)
}

type Updaters interface {
	UpdateCalculation(id, exp string) (Calculation, error)
}

type Deleters interface {
	DeleteCalculation(id string) error
}

type AllCalculationsGeter interface {
	GetAllCalculations() ([]Calculation, error)
}

type CalculationGeter interface {
	GetCalculation(id string) (Calculation, error)
}

// .
// .
// .

type calcService struct {
	repo CalculationRepository
}

func NewCalcService(r CalculationRepository) CalculationService {
	return &calcService{repo: r}
}

func (s *calcService) CalcExp(c *Calculation) error {
	return c.CalculatExpression()
}

// .
// .
// .

func (s *calcService) CreateCalculation(exp string) (Calculation, error) {
	// новый элемент
	calc := Calculation{
		ID:         uuid.NewString(),
		Expression: exp,
	}
	if err := s.CalcExp(&calc); err != nil {
		return Calculation{}, err
	}
	// сохраняем в БД
	if err := s.repo.CreateCalculation(calc); err != nil {
		return Calculation{}, err
	}
	return calc, nil
}

func (s *calcService) UpdateCalculation(id, exp string) (Calculation, error) {
	// находим в БД старую запись
	calc, err := s.repo.GetCalculation(id)
	if err != nil {
		return Calculation{}, err
	}

	// считаем новый результат
	calc.Expression = exp
	if err := s.CalcExp(&calc); err != nil {
		return Calculation{}, err
	}

	// обновляем запись в БД
	if err := s.repo.UpdateCalculation(calc); err != nil {
		return Calculation{}, err
	}

	return calc, nil
}

func (s *calcService) DeleteCalculation(id string) error {
	// мб добавить валидацию айдишника (дальше в репозиторий и бд не топаем)
	return s.repo.DeleteCalculation(id)
}

func (s *calcService) GetAllCalculations() ([]Calculation, error) {
	return s.repo.GetAllCalculations()
}

func (s *calcService) GetCalculation(id string) (Calculation, error) {
	return s.repo.GetCalculation(id)
}
