package dataset

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

const fileAlpha2 string = "alpha2.txt"
const directorySMS string = "dataset"

func initCountryAlpha2() (map[string]bool, error) {
	// Получаем текущий путь
	currentDir, err := os.Getwd()
	if err != nil {
		logrus.Errorf("Ошибка получения текущего пути: %s", err.Error())
		return nil, err
	}

	// Переходим на уровень вверх
	parentDir := filepath.Dir(currentDir)
	err = os.Chdir(parentDir)
	if err != nil {
		logrus.Errorf("Ошибка перехода на уровень вверх (aggregator): %s", err.Error())
		return nil, err
	}

	// Переходим в другую папку
	otherDir := filepath.Join(parentDir, directorySMS)
	err = os.Chdir(otherDir)
	if err != nil {
		logrus.Errorf("Ошибка перехода в другую папку (dataset): %s", err.Error())
		return nil, err
	}

	// Открываем файл
	file, err := os.Open(fileAlpha2)
	if err != nil {
		logrus.Errorf("Ошибка открытия файла (alpha2.txt): %s", err.Error())
		return nil, err
	}
	defer file.Close()

	// Создаем мапу для хранения кодов стран
	countriesMap := make(map[string]bool)

	// Создаем сканер для чтения файла построчно
	scanner := bufio.NewScanner(file)

	// Считываем каждую строку и разбиваем ее на коды стран
	for scanner.Scan() {
		line := scanner.Text()
		countries := strings.Fields(line)

		// Добавляем каждый код страны в мапу
		for _, country := range countries {
			countriesMap[country] = true
		}
	}

	// Проверяем наличие ошибок при сканировании файла
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return countriesMap, nil
}

func GetCountryMapAlpha2() (map[string]bool, error) {
	countriesMap, err := initCountryAlpha2()
	if err != nil {
		return nil, err
	}
	return countriesMap, nil
}
