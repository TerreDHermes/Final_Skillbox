package billing_handler

import (
	"aggregator"
	"fmt"
	"os"
)

const fileVoiceData string = "billing.data"

func getFileSimulatorEmailData() (*os.File, error) {
	// Открываем файл
	file, err := os.Open(fileVoiceData)
	if err != nil {
		fmt.Println("Ошибка открытия файла (billing.data):", err)
		return nil, err
	}
	return file, nil
}

func BillingHandler(dataBilling *SliceBilling) error {
	file, err := getFileSimulatorEmailData()
	if err != nil {
		return err
	}
	defer file.Close()

	// Чтение битовой маски из файла
	var mask string
	fmt.Fscanf(file, "%s", &mask)
	var number uint8
	for i := len(mask) - 1; i >= 0; i-- {
		if mask[i] == '1' {
			number += 1 << uint(len(mask)-1-i)
		}
	}
	billingData := aggregator.BillingData{
		CreateCustomer: number&(1<<0) != 0,
		Purchase:       number&(1<<1) != 0,
		Payout:         number&(1<<2) != 0,
		Recurring:      number&(1<<3) != 0,
		FraudControl:   number&(1<<4) != 0,
		CheckoutPage:   number&(1<<5) != 0,
	}
	dataBilling.AddBilling(billingData)
	return nil
}
