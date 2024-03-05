package support_handler

import (
	"aggregator"
)

type SliceSupport []aggregator.SupportData

func NewSliceSupport() SliceSupport {
	return make(SliceSupport, 0)
}

func (SupportSlice *SliceSupport) AddSupport(block aggregator.SupportData) {
	*SupportSlice = append(*SupportSlice, block)
}

func (SupportSlice SliceSupport) GetSupportLoadAndWaitTime() []int {
	// Средняя пропускная способность саппорта
	averageThroughput := 18
	// Количество специалистов
	numSpecialists := 7
	// Общее количество тикетов в час (предположим, что у есть поле с общим числом тикетов)
	totalTicketsPerHour := 0
	for _, supportData := range SupportSlice {
		totalTicketsPerHour += supportData.ActiveTickets
	}

	// Рассчитываем загруженность саппорта
	var supportLoad int
	switch {
	case totalTicketsPerHour <= averageThroughput*9:
		supportLoad = 1
	case totalTicketsPerHour <= averageThroughput*16:
		supportLoad = 2
	default:
		supportLoad = 3
	}

	// Рассчитываем среднее время ожидания ответа
	averageMinutesPerTicket := float64(60) / float64(averageThroughput*numSpecialists)
	waitTime := int(averageMinutesPerTicket * float64(totalTicketsPerHour))

	return []int{supportLoad, waitTime}
}
