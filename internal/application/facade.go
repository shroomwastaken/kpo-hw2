package application

import (
	"fmt"
	"pos/internal/domain"
)

type PosFacade struct {
	students    domain.StudentRepository
	rubricRepo  domain.RubricRepository
	evaluations domain.EvaluationRepository
	exporter    domain.ReportExporter
	factory     *domain.EvaluationFactory
	calcService *domain.GradeCalculationService
	seq         int
}

func NewPosFacade(
	students domain.StudentRepository,
	rubricRepo domain.RubricRepository,
	evaluations domain.EvaluationRepository,
	exporter domain.ReportExporter,
	factory *domain.EvaluationFactory,
	calcService *domain.GradeCalculationService,
) *PosFacade {
	return &PosFacade{
		students:    students,
		rubricRepo:  rubricRepo,
		evaluations: evaluations,
		exporter:    exporter,
		factory:     factory,
		calcService: calcService,
	}
}

func (f *PosFacade) nextID(prefix string) string {
	f.seq++
	return fmt.Sprintf("%s-%d", prefix, f.seq)
}

func (f *PosFacade) RegisterStudent(id, name, group string) (*domain.Student, error) {
	if id == "" {
		id = f.nextID("STU")
	}
	st, err := domain.NewStudent(domain.StudentID(id), name, group)
	if err != nil {
		return nil, err
	}
	if err := f.students.Save(st); err != nil {
		return nil, err
	}
	return st, nil
}

func (f *PosFacade) ListStudents() ([]*domain.Student, error) {
	return f.students.FindAll()
}

func (f *PosFacade) AddCriterion(id, title string, maxScore, weight float64) (*domain.Criterion, error) {
	if id == "" {
		id = f.nextID("CRT")
	}
	crit, err := domain.NewCriterion(domain.CriterionID(id), title, maxScore, weight)
	if err != nil {
		return nil, err
	}
	rubric, err := f.rubricRepo.Get()
	if err != nil {
		return nil, err
	}
	if err := rubric.AddCriterion(crit); err != nil {
		return nil, err
	}
	if err := f.rubricRepo.Save(rubric); err != nil {
		return nil, err
	}
	return crit, nil
}

func (f *PosFacade) ListCriteria() ([]*domain.Criterion, error) {
	rubric, err := f.rubricRepo.Get()
	if err != nil {
		return nil, err
	}
	return rubric.AllCriteria(), nil
}

func (f *PosFacade) AssignOrUpdateGrade(studentID, criterionID string, score float64, comment string) (*domain.Evaluation, error) {
	st, err := f.students.FindByID(domain.StudentID(studentID))
	if err != nil {
		return nil, err
	}
	rubric, err := f.rubricRepo.Get()
	if err != nil {
		return nil, err
	}

	existing, err := f.evaluations.FindByStudentAndCriterion(st.ID(), domain.CriterionID(criterionID))
	if err == nil && existing != nil {
		crit, err := rubric.FindCriterion(domain.CriterionID(criterionID))
		if err != nil {
			return nil, err
		}
		if err := existing.UpdateScore(score, crit, comment); err != nil {
			return nil, err
		}
		if err := f.evaluations.Save(existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	evalID := domain.EvaluationID(f.nextID("EVL"))
	eval, err := f.factory.Create(evalID, st, rubric, domain.CriterionID(criterionID), score, comment)
	if err != nil {
		return nil, err
	}
	if err := f.evaluations.Save(eval); err != nil {
		return nil, err
	}
	return eval, nil
}

func (f *PosFacade) RevokeGrade(evalID string) error {
	return f.evaluations.Delete(domain.EvaluationID(evalID))
}

func (f *PosFacade) ListStudentEvaluations(studentID string) ([]*domain.Evaluation, error) {
	return f.evaluations.FindByStudent(domain.StudentID(studentID))
}

func (f *PosFacade) SetGradingStrategy(strategy domain.GradingStrategy) {
	f.calcService.SetStrategy(strategy)
}

func (f *PosFacade) CurrentStrategyName() string {
	return f.calcService.StrategyName()
}

func (f *PosFacade) CalculateStudentFinal(studentID string) (domain.FinalGrade, error) {
	st, err := f.students.FindByID(domain.StudentID(studentID))
	if err != nil {
		return domain.FinalGrade{}, err
	}
	rubric, err := f.rubricRepo.Get()
	if err != nil {
		return domain.FinalGrade{}, err
	}
	evals, err := f.evaluations.FindByStudent(st.ID())
	if err != nil {
		return domain.FinalGrade{}, err
	}
	return f.calcService.CalculateForStudent(st, rubric, evals)
}

func (f *PosFacade) CalculateAllFinals() ([]domain.FinalGrade, error) {
	students, err := f.students.FindAll()
	if err != nil {
		return nil, err
	}
	rubric, err := f.rubricRepo.Get()
	if err != nil {
		return nil, err
	}
	var grades []domain.FinalGrade
	for _, st := range students {
		evals, err := f.evaluations.FindByStudent(st.ID())
		if err != nil {
			return nil, err
		}
		grade, err := f.calcService.CalculateForStudent(st, rubric, evals)
		if err != nil {
			return nil, err
		}
		grades = append(grades, grade)
	}
	return grades, nil
}

func (f *PosFacade) ExportFinals(filepath string) error {
	grades, err := f.CalculateAllFinals()
	if err != nil {
		return err
	}
	return f.exporter.Export(grades, filepath)
}
