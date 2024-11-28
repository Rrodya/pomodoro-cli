package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	workTime     = 25 // Время работы (минуты)
	breakTime    = 5  // Время перерыва (минуты)
	timerRunning = false
	isPaused     = false
	isWorking    = true
	remaining    time.Duration
)

// Форматирует время в виде MM:SS
func formatTime(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// Функция таймера
func startTimer(duration time.Duration, next func()) {
	remaining = duration
	timerRunning = true
	isPaused = false

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Очистка экрана и вывод только времени
	for remaining > 0 && timerRunning {
		if !isPaused {
			// Очистить строку и обновить ее
			fmt.Printf("\rОсталось: %s", formatTime(remaining))
			remaining -= 1 * time.Second
		}
		<-ticker.C
	}

	// Завершение работы таймера
	if timerRunning && remaining <= 0 {
		fmt.Printf("\nВремя %s закончилось!\n", map[bool]string{true: "работы", false: "перерыва"}[isWorking])
		next()
	}
}

// Меню выбора
func menu() {
	for {
		fmt.Println("\n=== Помодорро Таймер ===")
		fmt.Println("1. Установить время работы (текущее:", workTime, "мин.)")
		fmt.Println("2. Установить время перерыва (текущее:", breakTime, "мин.)")
		fmt.Println("3. Запустить таймер")
		fmt.Println("4. Остановить таймер")
		fmt.Println("5. Выход")
		fmt.Print("Выберите действие: ")

		input := getUserInput()
		switch input {
		case "1":
			fmt.Print("Введите время работы (в минутах): ")
			workTime = getUserInputInt()
			fmt.Println("Время работы установлено на", workTime, "минут.")
		case "2":
			fmt.Print("Введите время перерыва (в минутах): ")
			breakTime = getUserInputInt()
			fmt.Println("Время перерыва установлено на", breakTime, "минут.")
		case "3":
			if timerRunning {
				fmt.Println("Таймер уже запущен!")
			} else {
				isWorking = true
				runPomodoro()
			}
		case "4":
			stopTimer()
		case "5":
			fmt.Println("До свидания!")
			os.Exit(0)
		default:
			// Выводим ошибку на новой строке, чтобы она не затирала таймер
			fmt.Println("\nНеверный выбор. Попробуйте снова.")
		}
	}
}

// Запуск таймера помодорро
func runPomodoro() {
	go func() {
		startTimer(time.Duration(workTime)*time.Minute, func() {
			isWorking = false
			startTimer(time.Duration(breakTime)*time.Minute, func() {
				fmt.Println("Цикл завершён!")
				menu()
			})
		})
	}()

	// Управление таймером
	controlLoop()
}

// Управление таймером во время работы
func controlLoop() {
	for timerRunning {
		fmt.Print("\n[p] Пауза | [r] Продолжить | [s] Остановить: ")
		input := getUserInput()
		switch input {
		case "p":
			if timerRunning && !isPaused {
				isPaused = true
				fmt.Println("Таймер поставлен на паузу.")
			} else {
				fmt.Println("Таймер уже на паузе.")
			}
		case "r":
			if timerRunning && isPaused {
				isPaused = false
				fmt.Println("Таймер продолжен.")
			} else {
				fmt.Println("Таймер уже работает.")
			}
		case "s":
			stopTimer()
			menu()
			return
		default:
			// Ошибки выводим на новой строке
			fmt.Println("\nНеверная команда. Попробуйте снова.")
		}
	}
}

// Остановка таймера
func stopTimer() {
	if timerRunning {
		timerRunning = false
		isPaused = false
		fmt.Println("Таймер остановлен.")
	} else {
		fmt.Println("Таймер не запущен.")
	}
}

// Получение строки ввода пользователя
func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// Получение числового ввода пользователя
func getUserInputInt() int {
	input := getUserInput()
	value, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Ошибка ввода. Попробуйте снова.")
		return getUserInputInt()
	}
	return value
}

func main() {
	menu()
}
