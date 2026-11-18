package storage

import (
	"fmt"
	"pos/internal/domain"
	"sort"
	"sync"
)

type InMemoryStudentRepo struct {
	mu   sync.RWMutex
	data map[domain.StudentID]*domain.Student
}

func NewInMemoryStudentRepo() *InMemoryStudentRepo {
	return &InMemoryStudentRepo{data: make(map[domain.StudentID]*domain.Student)}
}

func (r *InMemoryStudentRepo) Save(s *domain.Student) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[s.ID()] = s
	return nil
}

func (r *InMemoryStudentRepo) FindByID(id domain.StudentID) (*domain.Student, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.data[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", domain.ErrStudentNotFound, id)
	}
	return s, nil
}

func (r *InMemoryStudentRepo) FindAll() ([]*domain.Student, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Student, 0, len(r.data))
	for _, s := range r.data {
		res = append(res, s)
	}
	sort.Slice(res, func(i, j int) bool { return res[i].ID() < res[j].ID() })
	return res, nil
}

type InMemoryRubricRepo struct {
	mu     sync.RWMutex
	rubric *domain.Rubric
}

func NewInMemoryRubricRepo() *InMemoryRubricRepo {
	return &InMemoryRubricRepo{rubric: domain.NewRubric()}
}

func (r *InMemoryRubricRepo) Get() (*domain.Rubric, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.rubric, nil
}

func (r *InMemoryRubricRepo) Save(rubric *domain.Rubric) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rubric = rubric
	return nil
}

type InMemoryEvaluationRepo struct {
	mu   sync.RWMutex
	data map[domain.EvaluationID]*domain.Evaluation
}

func NewInMemoryEvaluationRepo() *InMemoryEvaluationRepo {
	return &InMemoryEvaluationRepo{data: make(map[domain.EvaluationID]*domain.Evaluation)}
}

func (r *InMemoryEvaluationRepo) Save(e *domain.Evaluation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[e.ID()] = e
	return nil
}

func (r *InMemoryEvaluationRepo) FindByID(id domain.EvaluationID) (*domain.Evaluation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.data[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", domain.ErrEvaluationNotFound, id)
	}
	return e, nil
}

func (r *InMemoryEvaluationRepo) FindByStudent(studentID domain.StudentID) ([]*domain.Evaluation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*domain.Evaluation
	for _, e := range r.data {
		if e.StudentID() == studentID {
			res = append(res, e)
		}
	}
	return res, nil
}

func (r *InMemoryEvaluationRepo) FindByStudentAndCriterion(studentID domain.StudentID, criterionID domain.CriterionID) (*domain.Evaluation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, e := range r.data {
		if e.StudentID() == studentID && e.CriterionID() == criterionID {
			return e, nil
		}
	}
	return nil, domain.ErrEvaluationNotFound
}

func (r *InMemoryEvaluationRepo) Delete(id domain.EvaluationID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return fmt.Errorf("%w: %s", domain.ErrEvaluationNotFound, id)
	}
	delete(r.data, id)
	return nil
}

func (r *InMemoryEvaluationRepo) FindAll() ([]*domain.Evaluation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Evaluation, 0, len(r.data))
	for _, e := range r.data {
		res = append(res, e)
	}
	return res, nil
}
