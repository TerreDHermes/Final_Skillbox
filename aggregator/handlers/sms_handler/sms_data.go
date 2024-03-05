package sms_handler

import (
	"aggregator"
	"github.com/domodwyer/countrycodes"
	"github.com/sirupsen/logrus"
	"sort"
)

type SliceSms []aggregator.SMSDataMMS

func NewSliceSMS() SliceSms {
	return make(SliceSms, 0)
}

func (smsSlice *SliceSms) AddSMS(sms aggregator.SMSDataMMS) {
	*smsSlice = append(*smsSlice, sms)
}

func (smsSlice SliceSms) ProcessAndSortByProvider() SliceSms {
	processedSlice := make(SliceSms, len(smsSlice))
	copy(processedSlice, smsSlice)
	CountryMod(&processedSlice)
	// Сортируем по названию провайдера от A до Z
	sort.Slice(processedSlice, func(i, j int) bool {
		return processedSlice[i].Provider < processedSlice[j].Provider
	})
	return processedSlice
}

func (smsSlice SliceSms) ProcessAndSortByCountry() SliceSms {
	processedSlice := make(SliceSms, len(smsSlice))
	copy(processedSlice, smsSlice)
	CountryMod(&processedSlice)
	// Сортируем по стране от A до Z
	sort.Slice(processedSlice, func(i, j int) bool {
		return processedSlice[i].Country < processedSlice[j].Country
	})
	return processedSlice
}

func CountryMod(processedSlice *SliceSms) {
	var bl = "Saint Barthelemy"
	for i := range *processedSlice {
		name, err := countrycodes.ToName((*processedSlice)[i].Country)
		if err != nil {
			if (*processedSlice)[i].Country == "BL" {
				name = bl
			} else {
				logrus.Infof("Bad convert country: %s", (*processedSlice)[i].Country)
			}
		}
		(*processedSlice)[i].Country = name
	}
}
