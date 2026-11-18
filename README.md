# kpo hw2 (пос)

## сборка и запуск

```
$ go build ./cmd/pos

$ ./pos

# запуск демо
$ ./pos --demo
```

## карта ddd

* ubiquitous language:
  * `Student` - кого оценивают.
  * `Criterion` - за что ставят балл (`maxScore`, `weight`).
  * `Rubric` - набор актуальных критериев курса.
  * `Evaluation` - зафиксированный балл студента по конкретному критерию с комментарием и временем выставления.
  * `FinalGrade` - согласованный результат по студенту, вычисляемый из рубрикатора и оценок.
* entities / value objects:
  * entities: `Student`, `Criterion`, `Evaluation` - обладают уникальным идентификатором (`StudentID`, `CriterionID`, `EvaluationID`) и жизненным циклом (оценку можно обновить или отозвать).
  * value objects: `Score`, `Weight`, `FinalGrade`, типизированные ID - неизменяемы, не имеют собственной идентичности и валидируются в конструкторах (`NewScore`, `NewWeight`).
* инварианты:
  * балл не может быть меньше `0` или больше `criterion.MaxScore()`.
  * невозможно создать `Evaluation` по критерию, отсутствующему в `Rubric`.
  * итоговый балл не хранится как произвольное число, а всегда вычисляется согласованно из текущего состояния оценок и рубрикатора.
* repositories (порты): интерфейсы `StudentRepository`, `RubricRepository`, `EvaluationRepository` объявлены в `internal/domain/ports.go`.
* domain service: `GradeCalculationService` рассчитывает итоговую оценку по студенту, так как эта операция требует координации трёх понятий (`Student`, `Rubric`, `[]Evaluation`) и не принадлежит ни одной сущности единолично.

## архитектура

зависимости направлены к `domain`:

```
[console ui] -> [application] -> domain
                                   ^
[storage / csv exporter]------------
```

* `internal/domain` - чистое ядро без импортов инфраструктуры.
* `internal/application` - пользовательские сценарии (`PosFacade`). зависит только от `domain`.
* `internal/infrastructure` - внешние адаптеры: хранение в памяти (`storage`), экспорт в csv (`exporter`) и консольный интерфейс (`cli`).
* `cmd/pos/main.go` - di-контейнер, где создаются конкретные адаптеры и внедряются в сервисы через конструкторы.

## паттерны GoF

1. порождающий - factory (`EvaluationFactory` в `internal/domain/factory.go`):
   * обоснование: не даёт создать объект `Evaluation` в невалидном состоянии. проверяет наличие студента, наличие критерия в рубрикаторе и попадание `0 <= score <= max_score` в единой точке создания.
2. структурный - facade (`PosFacade` в `internal/application/facade.go`):
   * обоснование: изолирует консольный UI от деталей взаимодействия с тремя репозиториями, фабрикой и сервисом расчёта. консоль вызывает простые методы сценариев и отвечает только за ввод/вывод.
3. поведенческий - strategy (`GradingStrategy` в `internal/domain/service.go`):
   * обоснование: правила расчёта итога могут меняться (стандартная взвешенная шкала `WeightedTenScaleStrategy` или округление в пользу студента `RoundedUpGradingStrategy`). паттерн позволяет переключать алгоритм в рантайме без изменения `GradeCalculationService`.

## роль ии
ии не использовалось