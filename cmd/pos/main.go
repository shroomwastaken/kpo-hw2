package main

import (
	"os"
	"pos/internal/application"
	"pos/internal/domain"
	"pos/internal/infrastructure/cli"
	"pos/internal/infrastructure/exporter"
	"pos/internal/infrastructure/storage"
)

type Container struct {
	ConsoleUI *cli.ConsoleUI
}

func BuildContainer() *Container {
	studentRepo := storage.NewInMemoryStudentRepo()
	rubricRepo := storage.NewInMemoryRubricRepo()
	evalRepo := storage.NewInMemoryEvaluationRepo()
	csvExporter := exporter.NewCSVReportExporter()

	evalFactory := domain.NewEvaluationFactory()
	defaultStrategy := domain.WeightedTenScaleStrategy{}
	calcService := domain.NewGradeCalculationService(defaultStrategy)

	facade := application.NewPosFacade(
		studentRepo,
		rubricRepo,
		evalRepo,
		csvExporter,
		evalFactory,
		calcService,
	)

	consoleUI := cli.NewConsoleUI(facade)

	return &Container{
		ConsoleUI: consoleUI,
	}
}

func main() {
	container := BuildContainer()

	if len(os.Args) > 1 && os.Args[1] == "--demo" {
		container.ConsoleUI.RunDemoScenario()
		return
	}

	container.ConsoleUI.Run()
}
