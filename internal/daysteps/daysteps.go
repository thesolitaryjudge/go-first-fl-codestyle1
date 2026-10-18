package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage парсит строку с данными пакета активности и возвращает количество шагов,
// продолжительность и ошибку, если формат строки некорректен.
func parsePackage(data string) (int, time.Duration, error) {
	// 1. Разделить строку на слайс строк
	parts := strings.Split(data, ",")

	// 2. Проверить, чтобы длина слайса была равна 2
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат данных: ожидается 2 элемента через запятую")
	}

	// 3. Преобразовать первый элемент (количество шагов) в тип int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга шагов: %w", err)
	}

	// 4. Проверить: количество шагов должно быть больше 0
	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть больше 0")
	}

	// 5. Преобразовать второй элемент в time.Duration
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга продолжительности: %w", err)
	}

	// ИСПРАВЛЕНО: Добавлена проверка на нулевую или отрицательную продолжительность
	if duration <= 0 {
		return 0, 0, errors.New("продолжительность должна быть больше 0")
	}

	// 6. Если всё прошло без ошибок, вернуть результат
	return steps, duration, nil
}

// DayActionInfo формирует итоговую текстовую сводку об активности за день,
// включая количество шагов, пройденную дистанцию в километрах и потраченные калории.
func DayActionInfo(data string, weight, height float64) string {
	// 1. Получить данные о шагах и продолжительности
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("Ошибка:", err) // Вывести ошибку на экран
		return ""                   // Вернуть пустую строку
	}

	// 2. Проверить, чтобы количество шагов было больше 0
	if steps <= 0 {
		return ""
	}

	// 3. Вычислить дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// 4. Перевести дистанцию в километры
	distanceKm := distanceMeters / mInKm

	// 5. Вычислить количество калорий (вызов функции из пакета spentcalories)
	// ИСПРАВЛЕНО: теперь принимаем и калории, и ошибку
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println("Ошибка расчета калорий:", err)
		return ""
	}

	// 6. Сформировать и вернуть строку
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories)
}
