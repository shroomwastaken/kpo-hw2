package domain

type StudentRepository interface {
	Save(student *Student) error
	FindByID(id StudentID) (*Student, error)
	FindAll() ([]*Student, error)
}

type RubricRepository interface {
	Get() (*Rubric, error)
	Save(rubric *Rubric) error
}

type EvaluationRepository interface {
	Save(eval *Evaluation) error
	FindByID(id EvaluationID) (*Evaluation, error)
	FindByStudent(studentID StudentID) ([]*Evaluation, error)
	FindByStudentAndCriterion(studentID StudentID, criterionID CriterionID) (*Evaluation, error)
	Delete(id EvaluationID) error
	FindAll() ([]*Evaluation, error)
}

type ReportExporter interface {
	Export(grades []FinalGrade, destination string) error
}
