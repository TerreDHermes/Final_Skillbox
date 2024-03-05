package mms_handler

import (
	"aggregator"
	"github.com/domodwyer/countrycodes"
	"github.com/sirupsen/logrus"
	"sort"
)

type SliceMms []aggregator.SMSDataMMS

func NewSliceMMS() SliceMms {
	return make(SliceMms, 0)
}

func (mmsSlice *SliceMms) AddMMS(sms aggregator.SMSDataMMS) {
	*mmsSlice = append(*mmsSlice, sms)
}

func (mmsSlice SliceMms) ProcessAndSortByProvider() SliceMms {
	processedSlice := make(SliceMms, len(mmsSlice))
	copy(processedSlice, mmsSlice)
	CountryMod(&processedSlice)
	// Сортируем по названию провайдера от A до Z
	sort.Slice(processedSlice, func(i, j int) bool {
		return processedSlice[i].Provider < processedSlice[j].Provider
	})
	return processedSlice
}

func (mmsSlice SliceMms) ProcessAndSortByCountry() SliceMms {
	processedSlice := make(SliceMms, len(mmsSlice))
	copy(processedSlice, mmsSlice)
	CountryMod(&processedSlice)
	// Сортируем по названию провайдера от A до Z
	sort.Slice(processedSlice, func(i, j int) bool {
		return processedSlice[i].Country < processedSlice[j].Country
	})
	return processedSlice
}

func CountryMod(processedSlice *SliceMms) {
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
