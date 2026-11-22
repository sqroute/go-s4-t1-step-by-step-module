package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	spentcalories "github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	parts := strings.Split(data, ",")

	if len(parts) != 2 {
		return 0, 0, errors.New("incorrect data format")
	}

	stepCount, err := strconv.Atoi(parts[0])

	if err != nil {
		return 0, 0, err
	}

	if stepCount < 1 {
		return 0, 0, errors.New("one or more steps were expected")
	}

	walkDuration, err := time.ParseDuration(parts[1])

	if err != nil {
		return 0, 0, err
	}

	if walkDuration <= 0 {
		return 0, 0, errors.New("duration must be positive")
	}

	return stepCount, walkDuration, nil

}

func DayActionInfo(data string, weight, height float64) string {

	stepCount, walkDuration, err := parsePackage(data)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	if stepCount < 1 {
		return ""
	}

	distanceMeters := float64(stepCount) * stepLength
	distanceKilometers := distanceMeters / float64(mInKm)

	spentCals, err := spentcalories.WalkingSpentCalories(stepCount, weight, height, walkDuration)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", stepCount, distanceKilometers, spentCals)
}
