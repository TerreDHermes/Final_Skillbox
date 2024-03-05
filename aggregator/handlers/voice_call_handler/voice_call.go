package voice_call_handler

import (
	"aggregator"
	"encoding/csv"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
)

const fileVoiceData string = "voice.data"

func initProvider() (map[string]bool, error) {
	// Создаем мапу для хранения провайдеров
	providersMap := make(map[string]bool)
	providersMap["TransparentCalls"] = true
	providersMap["E-Voice"] = true
	providersMap["JustPhone"] = true
	return providersMap, nil
}

func getFileSimulatorVoiceCallData() (*os.File, error) {
	// Открываем файл
	file, err := os.Open(fileVoiceData)
	if err != nil {
		fmt.Println("Ошибка открытия файла (voice.data):", err)
		return nil, err
	}
	return file, nil
}

func VoiceCallHandler(dataVoiceCall *SliceVoiceCall, countriesMap map[string]bool) error {
	providersMap, err := initProvider()
	if err != nil {
		return err
	}
	file, err := getFileSimulatorVoiceCallData()
	if err != nil {
		return err
	}
	defer file.Close()
	// Создаем CSV Reader
	reader := csv.NewReader(file)
	reader.Comma = ';'
	var VoiceCallStruct aggregator.VoiceCallData

	for {
		// Читаем строку
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		// Проверяем, что в строке 4 поля
		if len(record) != 8 {
			continue
		}
		// проверка страны
		if countriesMap[record[0]] {
			VoiceCallStruct.Country = record[0]
		} else {
			continue
		}
		// проверка пропускной способности в процентах
		value, _ := strconv.Atoi(record[1])
		if value >= 0 && value <= 100 {
			VoiceCallStruct.Bandwidth = record[1]
		} else {
			continue
		}
		// вставка времени
		VoiceCallStruct.ResponseTime = record[2]
		// проверка провайдера
		if providersMap[record[3]] {
			VoiceCallStruct.Provider = record[3]
		} else {
			continue
		}
		// стабильность соединения float32
		floatNumber, err := strconv.ParseFloat(record[4], 32)
		if err != nil {
			logrus.Infof("Ошибка при преобразовании в float32", err)
			continue
		}
		VoiceCallStruct.ConnectionStability = float32(floatNumber)
		//  TTFB
		VoiceCallStruct.TTFB, err = strconv.Atoi(record[5])
		if err != nil {
			logrus.Infof("Ошибка при преобразовании в int, %s", err)
			continue
		}
		// чистота связи
		VoiceCallStruct.VoicePurity, err = strconv.Atoi(record[6])
		if err != nil {
			logrus.Infof("Ошибка при преобразовании в int, %s", err)
			continue
		}
		//медиана длительности звонка
		VoiceCallStruct.MedianOfCallsTime, err = strconv.Atoi(record[7])
		if err != nil {
			logrus.Infof("Ошибка при преобразовании в int, %s", err)
			continue
		}
		// добавляем строку в слайс (все проверки пройдены)
		(*dataVoiceCall).AddVoiceCall(VoiceCallStruct)
	}
	return nil
}
