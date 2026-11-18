package domain_test

import (
	"errors"
	"pos/internal/domain"
	"testing"
)

func TestEvaluationFactory_Invariants(t *testing.T) {
	student, _ := domain.NewStudent("S1", "Тестовый Студент", "БПИ231")
	rubric := domain.NewRubric()
	crit, _ := domain.NewCriterion("C1", "ДЗ №2", 10.0, 0.5)
	_ = rubric.AddCriterion(crit)

	factory := domain.NewEvaluationFactory()

	_, err := factory.Create("E1", student, rubric, "C1", 12.0, "Перебор")
	if !errors.Is(err, domain.ErrScoreOutOfBounds) {
		t.Fatalf("ожидалась ошибка ErrScoreOutOfBounds, получено: %v", err)
	}

	_, err = factory.Create("E2", student, rubric, "C1", -1.0, "Минус")
	if !errors.Is(err, domain.ErrScoreOutOfBounds) {
		t.Fatalf("ожидалась ошибка ErrScoreOutOfBounds, получено: %v", err)
	}

	_, err = factory.Create("E3", student, rubric, "C_MISSING", 8.0, "Нет критерия")
	if !errors.Is(err, domain.ErrCriterionNotFound) {
		t.Fatalf("ожидалась ошибка ErrCriterionNotFound, получено: %v", err)
	}

	ev, err := factory.Create("E4", student, rubric, "C1", 8.5, "Хорошо")
	if err != nil {
		t.Fatalf("неожиданная ошибка при валидной оценке: %v", err)
	}
	if ev.Score().Value() != 8.5 {
		t.Fatalf("ожидался балл 8.5, получено: %.2f", ev.Score().Value())
	}
}

func TestGradeCalculationService_Strategies(t *testing.T) {
	student, _ := domain.NewStudent("S1", "Иван Иванов", "БПИ231")
	rubric := domain.NewRubric()
	c1, _ := domain.NewCriterion("C1", "Семинары", 10.0, 0.4)
	c2, _ := domain.NewCriterion("C2", "Экзамен", 20.0, 0.6)
	_ = rubric.AddCriterion(c1)
	_ = rubric.AddCriterion(c2)

	factory := domain.NewEvaluationFactory()
	e1, _ := factory.Create("E1", student, rubric, "C1", 8.0, "")  // 8/10 * 10 * 0.4 = 3.2
	e2, _ := factory.Create("E2", student, rubric, "C2", 15.0, "") // 15/20 * 10 * 0.6 = 4.5

	svc := domain.NewGradeCalculationService(domain.WeightedTenScaleStrategy{})
	grade, err := svc.CalculateForStudent(student, rubric, []*domain.Evaluation{e1, e2})
	if err != nil {
		t.Fatalf("ошибка расчёта: %v", err)
	}
	if grade.RawTotal != 7.7 {
		t.Fatalf("ожидался итог 7.7, получено: %.2f", grade.RawTotal)
	}

	svc.SetStrategy(domain.RoundedUpGradingStrategy{})
	roundedGrade, err := svc.CalculateForStudent(student, rubric, []*domain.Evaluation{e1, e2})
	if err != nil {
		t.Fatalf("ошибка расчёта: %v", err)
	}
	if roundedGrade.RawTotal != 8.0 {
		t.Fatalf("ожидался округлённый итог 8.0, получено: %.2f", roundedGrade.RawTotal)
	}
}
