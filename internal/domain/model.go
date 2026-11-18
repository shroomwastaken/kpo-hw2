package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidStudent    = errors.New("некорректные данные студента")
	ErrInvalidCriterion  = errors.New("некорректные параметры критерия")
	ErrScoreOutOfBounds  = errors.New("балл выходит за пределы [0, max_score] критерия")
	ErrCriterionNotFound = errors.New("критерий отсутствует в рубрикаторе")
	ErrStudentNotFound   = errors.New("студент не найден")
	ErrEvaluationNotFound = errors.New("оценка не найдена")
	ErrEmptyRubric       = errors.New("рубрикатор пуст")
)

type StudentID string
type CriterionID string
type EvaluationID string

type Score struct {
	value float64
}

func NewScore(v float64) (Score, error) {
	if v < 0 {
		return Score{}, fmt.Errorf("%w: балл %.2f не может быть отрицательным", ErrScoreOutOfBounds, v)
	}
	return Score{value: v}, nil
}

func (s Score) Value() float64 { return s.value }

type Weight struct {
	value float64
}

func NewWeight(v float64) (Weight, error) {
	if v <= 0 || v > 1.0 {
		return Weight{}, fmt.Errorf("%w: вес %.2f должен быть в интервале (0, 1]", ErrInvalidCriterion, v)
	}
	return Weight{value: v}, nil
}

func (w Weight) Value() float64 { return w.value }

// FinalGrade — Value Object рассчитанного итога по студенту.
type FinalGrade struct {
	StudentID   StudentID
	StudentName string
	RawTotal    float64
	MaxPossible float64
	Details     []string
}

type Student struct {
	id    StudentID
	name  string
	group string
}

func NewStudent(id StudentID, name, group string) (*Student, error) {
	if strings.TrimSpace(string(id)) == "" || strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("%w: id и имя обязательны", ErrInvalidStudent)
	}
	return &Student{
		id:    id,
		name:  strings.TrimSpace(name),
		group: strings.TrimSpace(group),
	}, nil
}

func (s *Student) ID() StudentID  { return s.id }
func (s *Student) Name() string   { return s.name }
func (s *Student) Group() string  { return s.group }

type Criterion struct {
	id       CriterionID
	title    string
	maxScore Score
	weight   Weight
}

func NewCriterion(id CriterionID, title string, maxScore float64, weight float64) (*Criterion, error) {
	if strings.TrimSpace(string(id)) == "" || strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("%w: id и название критерия обязательны", ErrInvalidCriterion)
	}
	if maxScore <= 0 {
		return nil, fmt.Errorf("%w: max_score должен быть > 0", ErrInvalidCriterion)
	}
	ms, err := NewScore(maxScore)
	if err != nil {
		return nil, err
	}
	w, err := NewWeight(weight)
	if err != nil {
		return nil, err
	}
	return &Criterion{
		id:       id,
		title:    strings.TrimSpace(title),
		maxScore: ms,
		weight:   w,
	}, nil
}

func (c *Criterion) ID() CriterionID   { return c.id }
func (c *Criterion) Title() string     { return c.title }
func (c *Criterion) MaxScore() Score   { return c.maxScore }
func (c *Criterion) Weight() Weight    { return c.weight }

type Rubric struct {
	criteria map[CriterionID]*Criterion
}

func NewRubric() *Rubric {
	return &Rubric{criteria: make(map[CriterionID]*Criterion)}
}

func (r *Rubric) AddCriterion(c *Criterion) error {
	if c == nil {
		return ErrInvalidCriterion
	}
	if _, exists := r.criteria[c.ID()]; exists {
		return fmt.Errorf("%w: критерий %s уже существует", ErrInvalidCriterion, c.ID())
	}
	r.criteria[c.ID()] = c
	return nil
}

func (r *Rubric) FindCriterion(id CriterionID) (*Criterion, error) {
	c, ok := r.criteria[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrCriterionNotFound, id)
	}
	return c, nil
}

func (r *Rubric) AllCriteria() []*Criterion {
	res := make([]*Criterion, 0, len(r.criteria))
	for _, c := range r.criteria {
		res = append(res, c)
	}
	return res
}

type Evaluation struct {
	id          EvaluationID
	studentID   StudentID
	criterionID CriterionID
	score       Score
	comment     string
	gradedAt    time.Time
}

func (e *Evaluation) ID() EvaluationID           { return e.id }
func (e *Evaluation) StudentID() StudentID       { return e.studentID }
func (e *Evaluation) CriterionID() CriterionID   { return e.criterionID }
func (e *Evaluation) Score() Score               { return e.score }
func (e *Evaluation) Comment() string            { return e.comment }
func (e *Evaluation) GradedAt() time.Time        { return e.gradedAt }

func (e *Evaluation) UpdateScore(newScore float64, criterion *Criterion, comment string) error {
	if criterion == nil || criterion.ID() != e.criterionID {
		return ErrCriterionNotFound
	}
	sc, err := NewScore(newScore)
	if err != nil {
		return err
	}
	if sc.Value() > criterion.MaxScore().Value() {
		return fmt.Errorf("%w: %.2f > %.2f", ErrScoreOutOfBounds, sc.Value(), criterion.MaxScore().Value())
	}
	e.score = sc
	e.comment = comment
	e.gradedAt = time.Now()
	return nil
}
