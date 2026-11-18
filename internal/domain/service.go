package domain

import (
	"fmt"
	"math"
)

type GradingStrategy interface {
	Name() string
	Calculate(student *Student, rubric *Rubric, evaluations []*Evaluation) (FinalGrade, error)
}

type WeightedTenScaleStrategy struct{}

func (s WeightedTenScaleStrategy) Name() string {
	return "Взвешенная 10-балльная шкала (стандарт)"
}

func (s WeightedTenScaleStrategy) Calculate(student *Student, rubric *Rubric, evaluations []*Evaluation) (FinalGrade, error) {
	if student == nil {
		return FinalGrade{}, ErrStudentNotFound
	}
	criteria := rubric.AllCriteria()
	if len(criteria) == 0 {
		return FinalGrade{}, ErrEmptyRubric
	}

	evalByCrit := make(map[CriterionID]*Evaluation)
	for _, ev := range evaluations {
		if ev.StudentID() == student.ID() {
			evalByCrit[ev.CriterionID()] = ev
		}
	}

	var totalWeighted float64
	var weightSum float64
	var details []string

	for _, c := range criteria {
		weightSum += c.Weight().Value()
		var got float64
		if ev, ok := evalByCrit[c.ID()]; ok {
			got = ev.Score().Value()
		}
		normalized := (got / c.MaxScore().Value()) * 10.0 * c.Weight().Value()
		totalWeighted += normalized
		details = append(details, fmt.Sprintf("%s: %.2f/%.2f (вес %.2f -> +%.2f)",
			c.Title(), got, c.MaxScore().Value(), c.Weight().Value(), normalized))
	}

	if weightSum > 0 {
		totalWeighted = totalWeighted / weightSum
	}

	return FinalGrade{
		StudentID:   student.ID(),
		StudentName: student.Name(),
		RawTotal:    math.Round(totalWeighted*100) / 100,
		MaxPossible: 10.0,
		Details:     details,
	}, nil
}

type RoundedUpGradingStrategy struct {
	base WeightedTenScaleStrategy
}

func (s RoundedUpGradingStrategy) Name() string {
	return "Взвешенная 10-балльная шкала (округление в пользу студента)"
}

func (s RoundedUpGradingStrategy) Calculate(student *Student, rubric *Rubric, evaluations []*Evaluation) (FinalGrade, error) {
	grade, err := s.base.Calculate(student, rubric, evaluations)
	if err != nil {
		return FinalGrade{}, err
	}
	grade.RawTotal = math.Ceil(grade.RawTotal)
	return grade, nil
}

type GradeCalculationService struct {
	strategy GradingStrategy
}

func NewGradeCalculationService(strategy GradingStrategy) *GradeCalculationService {
	if strategy == nil {
		strategy = WeightedTenScaleStrategy{}
	}
	return &GradeCalculationService{strategy: strategy}
}

func (s *GradeCalculationService) SetStrategy(strategy GradingStrategy) {
	if strategy != nil {
		s.strategy = strategy
	}
}

func (s *GradeCalculationService) StrategyName() string {
	return s.strategy.Name()
}

func (s *GradeCalculationService) CalculateForStudent(
	student *Student,
	rubric *Rubric,
	evaluations []*Evaluation,
) (FinalGrade, error) {
	return s.strategy.Calculate(student, rubric, evaluations)
}
