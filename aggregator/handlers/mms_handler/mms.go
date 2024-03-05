package mms_handler

import (
	"aggregator"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

const UrlAddress = "http://127.0.0.1:8383/mms"

func initProvider() (map[string]bool, error) {
	// Создаем мапу для хранения кодов стран
	providersMap := make(map[string]bool)
	providersMap["Topolo"] = true
	providersMap["Rond"] = true
	providersMap["Kildy"] = true
	return providersMap, nil
}

func GetRequest() (*http.Response, error) {
	// Создаем новый запрос
	req, err := http.NewRequest("GET", UrlAddress, nil)
	if err != nil {
		return nil, err
	}
	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func JsonHandler(body []byte, dataMms *SliceMms, countriesMap map[string]bool, providersMap map[string]bool) error {
	// Unmarshal для декодирования JSON
	var responseData []aggregator.SMSDataMMS
	err := json.Unmarshal(body, &responseData)
	if err != nil {
		return err
	}
	// Проверка данных перед добавлением в слайс
	for _, entry := range responseData {
		if isValidEntry(entry, countriesMap, providersMap) {
			dataMms.AddMMS(entry)
		} else {
			logrus.Infof("Invalid entry detected in mms. Skipping.")
		}
	}
	return nil
}

// Функция для проверки валидности записи
func isValidEntry(entry aggregator.SMSDataMMS, countriesMap map[string]bool, providersMap map[string]bool) bool {
	// проверка страны
	if !countriesMap[entry.Country] {
		return false
	}
	// проверка пропускной способности в процентах
	value, _ := strconv.Atoi(entry.Bandwidth)
	if value < 0 || value > 100 {
		return false
	}
	// проверка провайдера
	if !providersMap[entry.Provider] {
		return false
	}
	return true
}

func MMSHandler(dataMms *SliceMms, countriesMap map[string]bool) error {
	providersMap, err := initProvider()
	if err != nil {
		return err
	}
	resp, err := GetRequest()
	if err != nil {
		return err
	} else if resp.Status != "200 OK" {
		logrus.Infof("Status not 200 in MMS response, %s", resp.Status)
	} else {
		// Читаем данные из ответа
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if err := JsonHandler(body, dataMms, countriesMap, providersMap); err != nil {
			return err
		}
	}
	defer resp.Body.Close()
	return nil
}
