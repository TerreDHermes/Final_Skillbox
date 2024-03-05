package sms_handler

import (
	"aggregator"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

const fileSmsData string = "sms.data"

func initProvider() (map[string]bool, error) {
	// Создаем мапу для хранения кодов стран
	providersMap := make(map[string]bool)
	providersMap["Topolo"] = true
	providersMap["Rond"] = true
	providersMap["Kildy"] = true
	return providersMap, nil
}

func getFileSimulatorSmsData() (*os.File, error) {
	file, err := os.Open(fileSmsData)
	if err != nil {
		fmt.Println("Ошибка открытия файла (sms.data):", err)
		return nil, err
	}
	return file, nil
}

func SMSHandler(dataSms *SliceSms, countriesMap map[string]bool) error {
	providersMap, err := initProvider()
	if err != nil {
		return err
	}
	file, err := getFileSimulatorSmsData()
	if err != nil {
		return err
	}
	defer file.Close()

	// Создаем CSV Reader
	reader := csv.NewReader(file)
	reader.Comma = ';'

	var smsStruct aggregator.SMSDataMMS

	// Читаем строки из файла
	for {
		// Читаем строку
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		// Проверяем, что в строке 4 поля
		if len(record) != 4 {
			continue
		}
		// проверка страны
		if countriesMap[record[0]] {
			smsStruct.Country = record[0]
		} else {
			continue
		}
		// проверка пропускной способности в процентах
		value, _ := strconv.Atoi(record[1])
		if value >= 0 && value <= 100 {
			smsStruct.Bandwidth = record[1]
		} else {
			break
		}
		// вставка времени
		smsStruct.ResponseTime = record[2]
		// проверка провайдера
		if providersMap[record[3]] {
			smsStruct.Provider = record[3]
			// добавляем строку в слайс (все проверки пройдены)
			(*dataSms).AddSMS(smsStruct)
		} else {
			break
		}
	}
	return nil
}
