package domain

import (
	"fmt"
	"time"
)

type EvaluationFactory struct{}

func NewEvaluationFactory() *EvaluationFactory {
	return &EvaluationFactory{}
}

func (f *EvaluationFactory) Create(
	id EvaluationID,
	student *Student,
	rubric *Rubric,
	criterionID CriterionID,
	rawScore float64,
	comment string,
) (*Evaluation, error) {
	if student == nil {
		return nil, ErrStudentNotFound
	}
	if rubric == nil {
		return nil, ErrEmptyRubric
	}
	criterion, err := rubric.FindCriterion(criterionID)
	if err != nil {
		return nil, err
	}
	score, err := NewScore(rawScore)
	if err != nil {
		return nil, err
	}
	if score.Value() > criterion.MaxScore().Value() {
		return nil, fmt.Errorf(
			"%w: передано %.2f при максимуме %.2f для критерия '%s'",
			ErrScoreOutOfBounds, score.Value(), criterion.MaxScore().Value(), criterion.Title(),
		)
	}

	return &Evaluation{
		id:          id,
		studentID:   student.ID(),
		criterionID: criterion.ID(),
		score:       score,
		comment:     comment,
		gradedAt:    time.Now(),
	}, nil
}
