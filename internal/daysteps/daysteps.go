package daysteps

import (
	"fmt"
	"log"
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

// parsePackage парсит строку вида "678,0h50m"
func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("некорректный формат: нужно 'шаги,длительность'")
	}

	stepsStr := parts[0]
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("не удалось преобразовать шаги: %w", err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть > 0")
	}

	durStr := parts[1]
	dur, err := time.ParseDuration(durStr)
	if err != nil {
		return 0, 0, fmt.Errorf("не удалось распарсить длительность: %w", err)
	}
	if dur <= 0 {
		return 0, 0, fmt.Errorf("длительность должна быть > 0")
	}

	return steps, dur, nil
}

// DayActionInfo возвращает инфу о дне в нужном формате
func DayActionInfo(data string, weight, height float64) string {
	steps, dur, err := parsePackage(data)
	if err != nil {
		log.Println("daysteps:", err)
		return ""
	}

	// расстояние в метрах
	distMeters := float64(steps) * stepLength
	// в км
	distKm := distMeters / float64(mInKm)

	// калории считаем через ходьбу
	cals, err := spentcalories.WalkingSpentCalories(steps, weight, height, dur)
	if err != nil {
		log.Println("daysteps: ошибка расчёта калорий:", err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps,
		distKm,
		cals,
	)
}
