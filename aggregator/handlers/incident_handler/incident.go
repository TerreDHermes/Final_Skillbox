package incident_handler

import (
	"aggregator"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const UrlAddress = "http://127.0.0.1:8383/accendent"

func initStatus() (map[string]bool, error) {
	// Создаем мапу для хранения кодов стран
	statusMap := make(map[string]bool)
	statusMap["active"] = true
	statusMap["closed"] = true
	return statusMap, nil
}

func JsonHandler(body []byte, dataIncident *SliceIncident) error {
	statusMap, err := initStatus()
	if err != nil {
		return err
	}
	// Unmarshal для декодирования JSON
	var responseData []aggregator.IncidentData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return err
	}
	// Проверка данных перед добавлением в слайс
	for _, entry := range responseData {
		// проверка статуса
		if statusMap[entry.Status] {
			dataIncident.AddIncident(entry)
		} else {
			logrus.Infof("Bad status: %s - %s", entry.Topic, entry.Status)
		}
	}
	return nil
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

func IncidentHandler(dataSupport *SliceIncident) error {
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
		if err := JsonHandler(body, dataSupport); err != nil {
			return err
		}
	}
	defer resp.Body.Close()
	return nil
}
