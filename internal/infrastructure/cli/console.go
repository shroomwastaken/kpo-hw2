package cli

import (
	"bufio"
	"fmt"
	"os"
	"pos/internal/application"
	"pos/internal/domain"
	"strconv"
	"strings"
)

type ConsoleUI struct {
	facade  *application.PosFacade
	scanner *bufio.Scanner
}

func NewConsoleUI(facade *application.PosFacade) *ConsoleUI {
	return &ConsoleUI{
		facade:  facade,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (ui *ConsoleUI) readLine(prompt string) string {
	fmt.Print(prompt)
	if ui.scanner.Scan() {
		return strings.TrimSpace(ui.scanner.Text())
	}
	return ""
}

func (ui *ConsoleUI) Run() {
	fmt.Println("=== Модуль «ПОС» — Преподавательская Оценка Студентов (НИУ ВШЭ) ===")
	for {
		fmt.Println("\nВыберите действие:")
		fmt.Println("1. Добавить студента")
		fmt.Println("2. Список студентов")
		fmt.Println("3. Добавить критерий в рубрикатор")
		fmt.Println("4. Показать рубрикатор")
		fmt.Println("5. Выставить / изменить оценку")
		fmt.Println("6. Посмотреть оценки студента")
		fmt.Println("7. Отменить (удалить) оценку")
		fmt.Println("8. Рассчитать итоги по всем студентам")
		fmt.Println("9. Сменить стратегию расчёта итога")
		fmt.Println("10. Экспортировать ведомость в CSV")
		fmt.Println("11. Запустить сквозной демонстрационный сценарий")
		fmt.Println("0. Выход")

		choice := ui.readLine("> ")
		switch choice {
		case "1":
			ui.handleAddStudent()
		case "2":
			ui.handleListStudents()
		case "3":
			ui.handleAddCriterion()
		case "4":
			ui.handleListCriteria()
		case "5":
			ui.handleAssignGrade()
		case "6":
			ui.handleListStudentGrades()
		case "7":
			ui.handleRevokeGrade()
		case "8":
			ui.handleCalculateFinals()
		case "9":
			ui.handleSwitchStrategy()
		case "10":
			ui.handleExport()
		case "11":
			ui.RunDemoScenario()
		case "0":
			fmt.Println("Завершение работы.")
			return
		default:
			fmt.Println("Неизвестная команда.")
		}
	}
}

func (ui *ConsoleUI) handleAddStudent() {
	id := ui.readLine("ID студента (Enter для автогенерации): ")
	name := ui.readLine("ФИО студента: ")
	group := ui.readLine("Группа: ")
	st, err := ui.facade.RegisterStudent(id, name, group)
	if err != nil {
		fmt.Printf("[Ошибка] %v\n", err)
		return
	}
	fmt.Printf("[OK] Студент добавлен: [%s] %s (%s)\n", st.ID(), st.Name(), st.Group())
}

func (ui *ConsoleUI) handleListStudents() {
	list, _ := ui.facade.ListStudents()
	if len(list) == 0 {
		fmt.Println("Список студентов пуст.")
		return
	}
	for _, s := range list {
		fmt.Printf("- [%s] %s (группа: %s)\n", s.ID(), s.Name(), s.Group())
	}
}

func (ui *ConsoleUI) handleAddCriterion() {
	id := ui.readLine("Код критерия (Enter для автогенерации): ")
	title := ui.readLine("Название критерия: ")
	maxStr := ui.readLine("Максимальный балл: ")
	weightStr := ui.readLine("Вес критерия: ")

	maxVal, err1 := strconv.ParseFloat(maxStr, 64)
	weightVal, err2 := strconv.ParseFloat(weightStr, 64)
	if err1 != nil || err2 != nil {
		fmt.Println("[Ошибка] Некорректный числовой формат.")
		return
	}
	c, err := ui.facade.AddCriterion(id, title, maxVal, weightVal)
	if err != nil {
		fmt.Printf("[Ошибка] %v\n", err)
		return
	}
	fmt.Printf("[OK] Критерий добавлен: [%s] %s (макс: %.2f, вес: %.2f)\n",
		c.ID(), c.Title(), c.MaxScore().Value(), c.Weight().Value())
}

func (ui *ConsoleUI) handleListCriteria() {
	list, _ := ui.facade.ListCriteria()
	if len(list) == 0 {
		fmt.Println("Рубрикатор пуст.")
		return
	}
	for _, c := range list {
		fmt.Printf("- [%s] %s | max: %.2f | вес: %.2f\n",
			c.ID(), c.Title(), c.MaxScore().Value(), c.Weight().Value())
	}
}

func (ui *ConsoleUI) handleAssignGrade() {
	stID := ui.readLine("ID студента: ")
	critID := ui.readLine("ID критерия: ")
	scoreStr := ui.readLine("Балл: ")
	comment := ui.readLine("Комментарий (необязательно): ")

	scoreVal, err := strconv.ParseFloat(scoreStr, 64)
	if err != nil {
		fmt.Println("[Ошибка] Балл должен быть числом.")
		return
	}
	ev, err := ui.facade.AssignOrUpdateGrade(stID, critID, scoreVal, comment)
	if err != nil {
		fmt.Printf("[Нарушение правила домена] %v\n", err)
		return
	}
	fmt.Printf("[OK] Оценка [%s] сохранена: студент=%s, критерий=%s, балл=%.2f\n",
		ev.ID(), ev.StudentID(), ev.CriterionID(), ev.Score().Value())
}

func (ui *ConsoleUI) handleListStudentGrades() {
	stID := ui.readLine("ID студента: ")
	evals, err := ui.facade.ListStudentEvaluations(stID)
	if err != nil {
		fmt.Printf("[Ошибка] %v\n", err)
		return
	}
	if len(evals) == 0 {
		fmt.Println("У студента пока нет оценок.")
		return
	}
	for _, e := range evals {
		fmt.Printf("- [%s] критерий=%s | балл=%.2f | комментарий: %s\n",
			e.ID(), e.CriterionID(), e.Score().Value(), e.Comment())
	}
}

func (ui *ConsoleUI) handleRevokeGrade() {
	evID := ui.readLine("ID оценки для удаления: ")
	if err := ui.facade.RevokeGrade(evID); err != nil {
		fmt.Printf("[Ошибка] %v\n", err)
		return
	}
	fmt.Println("[OK] Оценка удалена.")
}

func (ui *ConsoleUI) handleCalculateFinals() {
	fmt.Printf("Текущая стратегия: %s\n", ui.facade.CurrentStrategyName())
	grades, err := ui.facade.CalculateAllFinals()
	if err != nil {
		fmt.Printf("[Ошибка] %v\n", err)
		return
	}
	for _, g := range grades {
		fmt.Printf("Студент [%s] %s -> Итог: %.2f / %.2f\n",
			g.StudentID, g.StudentName, g.RawTotal, g.MaxPossible)
		for _, d := range g.Details {
			fmt.Printf("   • %s\n", d)
		}
	}
}

func (ui *ConsoleUI) handleSwitchStrategy() {
	fmt.Println("1. Стандартная взвешенная 10-балльная шкала")
	fmt.Println("2. Округление в пользу студента (ceil)")
	opt := ui.readLine("> ")
	if opt == "2" {
		ui.facade.SetGradingStrategy(domain.RoundedUpGradingStrategy{})
	} else {
		ui.facade.SetGradingStrategy(domain.WeightedTenScaleStrategy{})
	}
	fmt.Printf("[OK] Активна стратегия: %s\n", ui.facade.CurrentStrategyName())
}

func (ui *ConsoleUI) handleExport() {
	path := ui.readLine("Имя файла CSV (по умолчанию pos_grades.csv): ")
	if path == "" {
		path = "pos_grades.csv"
	}
	if err := ui.facade.ExportFinals(path); err != nil {
		fmt.Printf("[Ошибка] %v\n", err)
		return
	}
	fmt.Printf("[OK] Ведомость успешно экспортирована в %s\n", path)
}

func (ui *ConsoleUI) RunDemoScenario() {
	fmt.Println("\n--- ДЕМО ---")

	_, _ = ui.facade.RegisterStudent("S1", "Иванов Иван", "БПИ231")
	_, _ = ui.facade.RegisterStudent("S2", "Петрова Анна", "БПИ231")
	fmt.Println("1. Студенты S1 и S2 зарегистрированы.")

	_, _ = ui.facade.AddCriterion("C1", "Посещаемость лекций и семинаров", 10, 0.3)
	_, _ = ui.facade.AddCriterion("C2", "Соблюдение дедлайнов ДЗ", 10, 0.4)
	_, _ = ui.facade.AddCriterion("C3", "ИИ использовал, но понял сам", 10, 0.3)
	fmt.Println("2. Критерии C1, C2, C3 добавлены в рубрикатор.")

	_, _ = ui.facade.AssignOrUpdateGrade("S1", "C1", 9, "Активен на семинарах")
	_, _ = ui.facade.AssignOrUpdateGrade("S1", "C2", 8, "Один дедлайн впритык")
	_, _ = ui.facade.AssignOrUpdateGrade("S1", "C3", 10, "Отлично защитил архитектуру")

	_, _ = ui.facade.AssignOrUpdateGrade("S2", "C1", 10, "100% посещаемость")
	_, _ = ui.facade.AssignOrUpdateGrade("S2", "C2", 10, "Всё вовремя")
	_, _ = ui.facade.AssignOrUpdateGrade("S2", "C3", 9, "Хорошее понимание DDD")
	fmt.Println("3. Корректные оценки выставлены.")

	if _, err := ui.facade.AssignOrUpdateGrade("S1", "C1", 12, "Попытка поставить 12 из 10"); err != nil {
		fmt.Printf("4a. [Инвариант сработал] %v\n", err)
	}
	if _, err := ui.facade.AssignOrUpdateGrade("S1", "UNKNOWN", 5, "Несуществующий критерий"); err != nil {
		fmt.Printf("4b. [Инвариант сработал] %v\n", err)
	}

	fmt.Println("5. Итоговые оценки:")
	ui.handleCalculateFinals()

	demoFile := "pos_demo_report.csv"
	if err := ui.facade.ExportFinals(demoFile); err != nil {
		fmt.Printf("6. Ошибка экспорта: %v\n", err)
	} else {
		fmt.Printf("6. Итоговая ведомость сохранена в файл '%s'.\n", demoFile)
	}
	fmt.Println("--- КОНЕЦ ДЕМО ---")
}
