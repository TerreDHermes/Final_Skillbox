package email_handler

import (
	"aggregator"
	"encoding/csv"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
)

const fileVoiceData string = "email.data"

func initProvider() (map[string]bool, error) {
	// Создаем мапу для хранения провайдеров
	providersMap := make(map[string]bool)
	providersMap["Mail.ru"] = true
	providersMap["Yandex"] = true
	providersMap["Protonmail"] = true
	providersMap["GMX"] = true
	providersMap["RediffMail"] = true
	providersMap["Live"] = true
	providersMap["AOL"] = true
	providersMap["Comcast"] = true
	providersMap["Orange"] = true
	providersMap["MSN"] = true
	providersMap["Hotmail"] = true
	providersMap["Yahoo"] = true
	providersMap["Gmail"] = true
	return providersMap, nil
}

func getFileSimulatorEmailData() (*os.File, error) {
	// Открываем файл
	file, err := os.Open(fileVoiceData)
	if err != nil {
		logrus.Infof("Ошибка открытия файла (email.data): %s", err)
		return nil, err
	}
	return file, nil
}

func EmailHandler(dataEmail *SliceEmail, countriesMap map[string]bool) error {
	providersMap, err := initProvider()
	if err != nil {
		return err
	}
	file, err := getFileSimulatorEmailData()
	if err != nil {
		return err
	}
	defer file.Close()

	// Создаем CSV Reader
	reader := csv.NewReader(file)
	reader.Comma = ';'
	var EmailStruct aggregator.EmailData

	for {
		// Читаем строку
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		// Проверяем, что в строке 4 поля
		if len(record) != 3 {
			continue
		}
		// проверка страны
		if countriesMap[record[0]] {
			EmailStruct.Country = record[0]
		} else {
			continue
		}
		// проверка провайдера
		if providersMap[record[1]] {
			EmailStruct.Provider = record[1]
		} else {
			continue
		}
		//среднее время доставки писем в ms
		EmailStruct.DeliveryTime, err = strconv.Atoi(record[2])
		if err != nil {
			logrus.Infof("Ошибка при преобразовании в int, %s", err)
			continue
		}
		// добавляем строку в слайс (все проверки пройдены)
		(*dataEmail).AddEmail(EmailStruct)
	}
	return nil
}
