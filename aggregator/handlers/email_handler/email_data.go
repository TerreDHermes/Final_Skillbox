package email_handler

import (
	"aggregator"
)

type SliceEmail []aggregator.EmailData

func NewSliceEmail() SliceEmail {
	return make(SliceEmail, 0)
}

func (EmailSlice *SliceEmail) AddEmail(sms aggregator.EmailData) {
	*EmailSlice = append(*EmailSlice, sms)
}

// Метод для сортировки провайдеров по времени доставки в каждой стране
func (EmailSlice SliceEmail) SortByDeliveryTime() map[string][][]aggregator.EmailData {
	// Создаем карту для хранения отсортированных данных
	sortedData := make(map[string][][]aggregator.EmailData)

	// Создаем карту для хранения временных данных перед сортировкой
	tempData := make(map[string][]aggregator.EmailData)

	// Заполняем временные данные
	for _, email := range EmailSlice {
		tempData[email.Country] = append(tempData[email.Country], email)
	}

	// Сортируем данные по времени доставки
	for country, emails := range tempData {
		// Сортируем по возрастанию времени доставки
		sortByDeliveryTime(emails)

		// Разбиваем на 3 быстрых и 3 медленных провайдера
		fastProviders := emails[:3]
		slowProviders := emails[len(emails)-3:]

		// Записываем в карту
		sortedData[country] = [][]aggregator.EmailData{fastProviders, slowProviders}
	}

	return sortedData
}

// Вспомогательная функция для сортировки по времени доставки
func sortByDeliveryTime(emails []aggregator.EmailData) {
	n := len(emails)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if emails[j].DeliveryTime > emails[j+1].DeliveryTime {
				// меняем местами элементы, если текущий больше следующего
				emails[j], emails[j+1] = emails[j+1], emails[j]
			}
		}
	}
}
