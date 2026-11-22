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
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {

	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, errors.New("incorrect data format")
	}

	stepCount, err := strconv.Atoi(parts[0])

	if err != nil {
		return 0, "", 0, err
	}

	if stepCount <= 0 {
		return 0, "", 0, errors.New("steps must be positive")
	}

	walkDuration, err := time.ParseDuration(parts[2])

	if err != nil {
		return 0, "", 0, err
	}

	if walkDuration <= 0 {
		return 0, "", 0, errors.New("duration must be positive")
	}

	return stepCount, parts[1], walkDuration, nil

}

func distance(steps int, height float64) float64 {

	stepLength := stepLengthCoefficient * height
	distanceMeters := float64(steps) * stepLength
	distanceKilometers := distanceMeters / float64(mInKm)

	return distanceKilometers

}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		return 0
	}

	distanceKilometers := distance(steps, height)
	averageSpeed := distanceKilometers / duration.Hours()

	return averageSpeed

}

func TrainingInfo(data string, weight, height float64) (string, error) {

	var distanceInKilometers, averageSpeed, calories float64
	steps, trainType, duration, err := parseTraining(data)

	if err != nil {
		log.Println(err)
		// ТЗ требует, чтобы я лишь залогировал ошибку, хотя если не завершить функцию,
		// то битые данные пойдут дальше по логике и все равно приведут к ошибке.
		// Не стал отходить от ТЗ, предполагая, что это намеренное требование.
	}

	switch trainType {
	case "Бег":
		distanceInKilometers = distance(steps, height)
		averageSpeed = meanSpeed(steps, height, duration)
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		distanceInKilometers = distance(steps, height)
		averageSpeed = meanSpeed(steps, height, duration)
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	message := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		trainType, duration.Hours(), distanceInKilometers, averageSpeed, calories)

	return message, nil
}

func calculateSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		return 0, errors.New("steps cannot be zero or negative")
	}

	if weight <= 0 {
		return 0, errors.New("weight cannot be zero or negative")
	}

	if height <= 0 {
		return 0, errors.New("height cannot be zero or negative")
	}

	if duration <= 0 {
		return 0, errors.New("duration cannot be zero or negative")
	}

	averageSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	calories := (weight * averageSpeed * durationInMinutes) / minInH

	return calories, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	spentCalories, err := calculateSpentCalories(steps, weight, height, duration)

	if err != nil {
		return 0, err
	}

	return spentCalories, nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	spentCalories, err := calculateSpentCalories(steps, weight, height, duration)

	if err != nil {
		return 0, err
	}

	return spentCalories * walkingCaloriesCoefficient, nil

}
