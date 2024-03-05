package billing_handler

import (
	"aggregator"
)

type SliceBilling []aggregator.BillingData

func NewSliceBilling() SliceBilling {
	return make(SliceBilling, 0)
}

func (BillingSlice *SliceBilling) AddBilling(sms aggregator.BillingData) {
	*BillingSlice = append(*BillingSlice, sms)
}
