package support_handler

import (
	"aggregator"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const UrlAddress = "http://127.0.0.1:8383/support"

func JsonHandler(body []byte, dataSupport *SliceSupport) error {
	// Unmarshal для декодирования JSON
	var responseData []aggregator.SupportData
	err := json.Unmarshal(body, &responseData)
	if err != nil {
		return err
	}
	// Проверка данных перед добавлением в слайс
	for _, entry := range responseData {
		dataSupport.AddSupport(entry)
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

func SupportHandler(dataSupport *SliceSupport) error {
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
