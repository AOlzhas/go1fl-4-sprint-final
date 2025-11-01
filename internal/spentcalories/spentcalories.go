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
	lenStep                    = 0.65 // можно не использовать, но оставим
	mInKm                      = 1000 // количество метров в километре
	minInH                     = 60   // количество минут в часе
	stepLengthCoefficient      = 0.45 // коэффициент длины шага от роста
	walkingCaloriesCoefficient = 0.5  // поправка для ходьбы
)

// parseTraining парсит строку формата "3456,Ходьба,3h00m"
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("некорректный формат данных тренировки: нужно 'шаги,тип,длительность'")
	}

	// шаги
	stepsStr := parts[0]
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось преобразовать шаги: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть > 0")
	}

	// тип активности
	activity := parts[1]

	// длительность
	durStr := parts[2]
	dur, err := time.ParseDuration(durStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось распарсить длительность: %w", err)
	}
	// тут у тебя как раз и падал тест — он ждал ошибку на 0h0m и -1h30m
	if dur <= 0 {
		return 0, "", 0, fmt.Errorf("длительность должна быть > 0")
	}

	return steps, activity, dur, nil
}

// distance считает дистанцию в КИЛОМЕТРАХ по шагам и росту
func distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 {
		return 0
	}
	stepLen := height * stepLengthCoefficient // метры
	distMeters := float64(steps) * stepLen
	return distMeters / float64(mInKm)
}

// meanSpeed возвращает среднюю скорость в км/ч
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distKm := distance(steps, height)
	hours := duration.Hours()
	if hours == 0 {
		return 0
	}
	return distKm / hours
}

// RunningSpentCalories считает калории для бега
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("шаги должны быть > 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть > 0")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть > 0")
	}
	if duration <= 0 {
		return 0, errors.New("длительность должна быть > 0")
	}

	speed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()

	cals := (weight * speed * minutes) / float64(minInH)
	return cals, nil
}

// WalkingSpentCalories считает калории для ходьбы
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("шаги должны быть > 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть > 0")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть > 0")
	}
	if duration <= 0 {
		return 0, errors.New("длительность должна быть > 0")
	}

	speed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()

	base := (weight * speed * minutes) / float64(minInH)
	return base * walkingCaloriesCoefficient, nil
}

// TrainingInfo формирует человекочитаемый отчёт по тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, dur, err := parseTraining(data)
	if err != nil {
		log.Println("spentcalories:", err)
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, dur)
	hours := dur.Hours()

	var cals float64

	switch strings.ToLower(strings.TrimSpace(activity)) {
	case "бег", "run", "running":
		cals, err = RunningSpentCalories(steps, weight, height, dur)
		if err != nil {
			log.Println("spentcalories:", err)
			return "", err
		}
		activity = "Бег"
	case "ходьба", "walk", "walking":
		cals, err = WalkingSpentCalories(steps, weight, height, dur)
		if err != nil {
			log.Println("spentcalories:", err)
			return "", err
		}
		activity = "Ходьба"
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activity)
	}

	// тест у тебя ругался, что нет \n в конце
	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity,
		hours,
		dist,
		speed,
		cals,
	), nil
}
