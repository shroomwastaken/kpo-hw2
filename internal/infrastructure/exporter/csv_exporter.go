package exporter

import (
	"encoding/csv"
	"fmt"
	"os"
	"pos/internal/domain"
	"strings"
)

type CSVReportExporter struct{}

func NewCSVReportExporter() *CSVReportExporter {
	return &CSVReportExporter{}
}

func (e *CSVReportExporter) Export(grades []domain.FinalGrade, destination string) error {
	file, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("не удалось создать файл %s: %w", destination, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"StudentID", "StudentName", "FinalGrade", "MaxPossible", "Breakdown"}); err != nil {
		return err
	}

	for _, g := range grades {
		row := []string{
			string(g.StudentID),
			g.StudentName,
			fmt.Sprintf("%.2f", g.RawTotal),
			fmt.Sprintf("%.2f", g.MaxPossible),
			strings.Join(g.Details, " | "),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return writer.Error()
}
