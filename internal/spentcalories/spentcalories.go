package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе.
)

// parseTraining парсит строку с данными о тренировке и возвращает количество шагов,
// тип активности, продолжительность занятия и ошибку, если формат некорректен.
func parseTraining(data string) (int, string, time.Duration, error) {
	// 1. Разделить строку на слайс строк.
	parts := strings.Split(data, ",")

	// 2. Проверить, чтобы длина слайса была равна 3.
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных: ожидается 3 элемента через запятую")
	}

	// 3. Преобразовать первый элемент слайса (количество шагов) в тип int.
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга шагов: %w", err)
	}

	// Проверить, чтобы количество шагов было больше 0.
	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше 0")
	}

	activityType := parts[1]

	// 4. Преобразовать третий элемент слайса в time.Duration.
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга продолжительности: %w", err)
	}

	// Проверить, чтобы продолжительность была больше 0.
	if duration <= 0 {
		return 0, "", 0, errors.New("продолжительность тренировки должна быть больше 0")
	}

	// 5. Если всё прошло без ошибок, верните значения.
	return steps, activityType, duration, nil
}

// distance рассчитывает дистанцию в километрах на основе количества шагов и роста человека.
func distance(steps int, height float64) float64 {
	// 1. Рассчитать длину шага.
	stepLength := height * stepLengthCoefficient

	// 2. Умножить пройденное количество шагов на длину шага.
	distanceMeters := float64(steps) * stepLength

	// 3. Разделите полученное значение на число метров в километре mInKm.
	return distanceMeters / float64(mInKm)
}

// meanSpeed рассчитывает среднюю скорость движения в км/ч за время тренировки.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// 1. Проверить, что продолжительность duration больше 0.
	if duration <= 0 {
		return 0
	}

	// 2. Вычислить дистанцию с помощью distance().
	dist := distance(steps, height)

	// 3. Вычислить и вернуть среднюю скорость (дистанция / часы).
	return dist / duration.Hours()
}

// RunningSpentCalories рассчитывает калории, потраченные при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// 1. Проверить входные параметры на корректность.
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные параметры для расчета калорий бега")
	}

	// 2. Рассчитать среднюю скорость с помощью meanSpeed().
	speed := meanSpeed(steps, height, duration)

	// 3. Переведите продолжительность в минуты с помощью функции из пакета time.
	durationInMinutes := duration.Minutes()

	// 4. Рассчитать и вернуть количество калорий по формуле из ТЗ:
	// (weight * meanSpeed * durationInMinutes) / minInH
	calories := (weight * speed * durationInMinutes) / float64(minInH)

	return calories, nil
}

// WalkingSpentCalories рассчитывает калории, потраченные при ходьбе.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные параметры для расчета калорий ходьбы")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	// Базовый расчет калорий, умноженный на корректирующий коэффициент.
	baseCalories := (weight * speed * durationInMinutes) / float64(minInH)
	calories := baseCalories * walkingCaloriesCoefficient

	return calories, nil
}

// TrainingInfo собирает полную информацию о тренировке и возвращает её в виде строки.
func TrainingInfo(data string, weight, height float64) (string, error) {
	// 1. Получить значения из строки данных с помощью функции parseTraining(),
	// обработать возможные ошибки и вывести их в лог с помощью log.Println(err).
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64
	var formattedActivityType string

	// 2. Проверить вид тренировки с помощью switch.
	switch strings.ToLower(activityType) {
	case "бег":
		formattedActivityType = "Бег"
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "ходьба":
		formattedActivityType = "Ходьба"
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		// Если передан неизвестный тип тренировки, вернуть ошибку с текстом "неизвестный тип тренировки".
		return "", errors.New("неизвестный тип тренировки")
	}

	// Вычисляем метрики
	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	durationInHours := duration.Hours()

	// 3. Формируем и возвращаем результирующую строку строго по шаблону.
	// Внимание: проверяем точное совпадение букв и точек.
	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		formattedActivityType, durationInHours, dist, speed, calories,
	)

	return result, nil
}
