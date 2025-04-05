package calculationservice

import "gorm.io/gorm"

// CGUD
type CalculationRepository interface {
	CreatersRep
	GetersRep
	UpdatersRep
	DeletersRep
}

type GetersRep interface {
	AllCalculationsGeterRep
	CalculationGeterRep
}

type CreatersRep interface {
	CreateCalculation(Calculation) error
}

type UpdatersRep interface {
	UpdateCalculation(Calculation) error
}

type DeletersRep interface {
	DeleteCalculation(string) error
}

type AllCalculationsGeterRep interface {
	GetAllCalculations() ([]Calculation, error)
}

type CalculationGeterRep interface {
	GetCalculation(string) (Calculation, error)
}

// .
// .

type calcRepository struct {
	db *gorm.DB
}

func NewCalculationRepository(db *gorm.DB) CalculationRepository {
	return &calcRepository{db: db}
}

// .
// реализация методов
// .

func (r *calcRepository) CreateCalculation(calc Calculation) error {
	return r.db.Create(&calc).Error
}

func (r *calcRepository) GetAllCalculations() (calcs []Calculation, err error) {
	err = r.db.Find(&calcs).Error
	return
}

func (r *calcRepository) GetCalculation(id string) (Calculation, error) {
	var calc = Calculation{ID: id}
	err := r.db.First(&calc).Error
	return calc, err
}

func (r *calcRepository) UpdateCalculation(calc Calculation) error {
	return r.db.Save(calc).Error
}

func (r *calcRepository) DeleteCalculation(id string) error {
	return r.db.Delete(&Calculation{}, "id = ?", id).Error
}
