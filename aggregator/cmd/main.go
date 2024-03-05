package main

import (
	"aggregator"
	"aggregator/dataset"
	"aggregator/handlers/billing_handler"
	"aggregator/handlers/email_handler"
	"aggregator/handlers/incident_handler"
	"aggregator/handlers/mms_handler"
	"aggregator/handlers/sms_handler"
	"aggregator/handlers/support_handler"
	"aggregator/handlers/voice_call_handler"
	"aggregator/web"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

const (
	port                      = "8282"
	updatePeriod              = 30 * time.Second
	maxWaitPeriod             = 10 * time.Second
	directorySimulator string = "simulator"
)

func InitDir() error {
	// Получаем текущий путь
	currentDir, err := os.Getwd()
	if err != nil {
		logrus.Errorf("Ошибка получения текущего пути: %s", err.Error())
		return err
	}
	// Переходим на уровень вверх
	parentDir := filepath.Dir(currentDir)
	err = os.Chdir(parentDir)
	if err != nil {
		logrus.Errorf("Ошибка перехода на уровень вверх (aggregator): %s", err.Error())
		return err
	}
	// Переходим на уровень вверх
	parentDir = filepath.Dir(parentDir)
	err = os.Chdir(parentDir)
	if err != nil {
		logrus.Errorf("Ошибка перехода на уровень вверх (final): %s", err.Error())
		return err
	}
	// Переходим в другую папку
	otherDir := filepath.Join(parentDir, directorySimulator)
	err = os.Chdir(otherDir)
	if err != nil {
		logrus.Errorf("Ошибка перехода в другую папку (simulator): %s", err.Error())
		return err
	}
	return nil
}

func getData() (sms_handler.SliceSms, mms_handler.SliceMms, voice_call_handler.SliceVoiceCall, email_handler.SliceEmail, billing_handler.SliceBilling, support_handler.SliceSupport, incident_handler.SliceIncident) {
	dataSms := sms_handler.NewSliceSMS()
	dataMms := mms_handler.NewSliceMMS()
	dataVoiceCall := voice_call_handler.NewSliceVoiceCall()
	dataEmail := email_handler.NewSliceEmail()
	dataBilling := billing_handler.NewSliceBilling()
	dataSupport := support_handler.NewSliceSupport()
	dataIncident := incident_handler.NewSliceIncident()
	return dataSms, dataMms, dataVoiceCall, dataEmail, dataBilling, dataSupport, dataIncident
}

func FillStatusStruct(statusStruct *aggregator.StatusHandler, message string) {
	(*statusStruct).Mu.Lock()
	(*statusStruct).Status = false
	(*statusStruct).Error = (*statusStruct).Error + "	" + message
	(*statusStruct).Mu.Unlock()
}

func startHandlers(statusStruct *aggregator.StatusHandler, wg *sync.WaitGroup, countriesMap map[string]bool, dataSms *sms_handler.SliceSms, dataMms *mms_handler.SliceMms, dataVoiceCall *voice_call_handler.SliceVoiceCall, dataEmail *email_handler.SliceEmail, dataBilling *billing_handler.SliceBilling, dataSupport *support_handler.SliceSupport, dataIncident *incident_handler.SliceIncident) {
	if err := sms_handler.SMSHandler(dataSms, countriesMap); err != nil {
		logrus.Infof("Error in SMS_Handler: %s", err.Error())
		FillStatusStruct(statusStruct, "Error in SMS_Handler")
	}
	go func() {
		defer wg.Done()
		if err := mms_handler.MMSHandler(dataMms, countriesMap); err != nil {
			logrus.Infof("Error in MMS_Handler: %s", err.Error())
			FillStatusStruct(statusStruct, "Error in MMS_Handler")
		}
	}()
	go func() {
		defer wg.Done()
		if err := voice_call_handler.VoiceCallHandler(dataVoiceCall, countriesMap); err != nil {
			logrus.Infof("Error in Voice_Call_Handler: %s", err.Error())
			FillStatusStruct(statusStruct, "Error in Voice_Call_Handler")
		}
	}()
	go func() {
		defer wg.Done()
		if err := email_handler.EmailHandler(dataEmail, countriesMap); err != nil {
			logrus.Infof("Error in Email_Handler: %s", err.Error())
			FillStatusStruct(statusStruct, "Error in Email_Handler")
		}
	}()
	go func() {
		defer wg.Done()
		if err := billing_handler.BillingHandler(dataBilling); err != nil {
			logrus.Infof("Error in Billing_Handler: %s", err.Error())
			FillStatusStruct(statusStruct, "Error in Billing_Handler")
		}
	}()
	go func() {
		defer wg.Done()
		if err := support_handler.SupportHandler(dataSupport); err != nil {
			logrus.Infof("Error in Support_Handler: %s", err.Error())
			FillStatusStruct(statusStruct, "Error in Support_Handler")
		}
	}()
	go func() {
		defer wg.Done()
		if err := incident_handler.IncidentHandler(dataIncident); err != nil {
			logrus.Infof("Error in Incident_Handler: %s", err.Error())
			FillStatusStruct(statusStruct, "Error in Incident_Handler")
		}
	}()
}

func CreateResponse(statusStruct *aggregator.StatusHandler, dataSms *sms_handler.SliceSms, dataMms *mms_handler.SliceMms, dataVoiceCall *voice_call_handler.SliceVoiceCall, dataEmail *email_handler.SliceEmail, dataBilling *billing_handler.SliceBilling, dataSupport *support_handler.SliceSupport, dataIncident *incident_handler.SliceIncident) ([]byte, error) {
	var data aggregator.ResultSetT
	if (*statusStruct).Status {
		sms := [][]aggregator.SMSDataMMS{(*dataSms).ProcessAndSortByProvider(), (*dataSms).ProcessAndSortByCountry()}
		mms := [][]aggregator.SMSDataMMS{(*dataMms).ProcessAndSortByProvider(), (*dataMms).ProcessAndSortByCountry()}
		voiceCall := *dataVoiceCall
		email := dataEmail.SortByDeliveryTime()
		billing := aggregator.BillingData{}
		if len(*dataBilling) > 0 {
			billing = (*dataBilling)[0]
		}
		support := (*dataSupport).GetSupportLoadAndWaitTime()
		incidents := (*dataIncident).SortByStatus()

		data = aggregator.ResultSetT{
			SMS:       sms,
			MMS:       mms,
			VoiceCall: voiceCall,
			Email:     email,
			Billing:   billing,
			Support:   support,
			Incidents: incidents,
		}
	}

	result := aggregator.ResultT{
		Status: (*statusStruct).Status,
		Data:   data,
		Error:  (*statusStruct).Error,
	}

	JsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		logrus.Errorf("Error json.MarshalIndent: %s", err.Error())
		return nil, err
	}
	return JsonData, nil
}

func updateDataPeriodically(statusStruct *aggregator.StatusHandler, handler *web.Handler, countriesMap map[string]bool, dataSms *sms_handler.SliceSms, dataMms *mms_handler.SliceMms, dataVoiceCall *voice_call_handler.SliceVoiceCall, dataEmail *email_handler.SliceEmail, dataBilling *billing_handler.SliceBilling, dataSupport *support_handler.SliceSupport, dataIncident *incident_handler.SliceIncident) {
	var wg2 sync.WaitGroup

	timer := time.NewTicker(updatePeriod)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			// Обновляем данные
			(*statusStruct).Status = true
			(*statusStruct).Error = ""
			*dataSms, *dataMms, *dataVoiceCall, *dataEmail, *dataBilling, *dataSupport, *dataIncident = nil, nil, nil, nil, nil, nil, nil
			wg2.Add(6)
			startHandlers(statusStruct, &wg2, countriesMap, dataSms, dataMms, dataVoiceCall, dataEmail, dataBilling, dataSupport, dataIncident)
			wg2.Wait()

			JsonData, err := CreateResponse(statusStruct, dataSms, dataMms, dataVoiceCall, dataEmail, dataBilling, dataSupport, dataIncident)
			if err != nil {
				logrus.Errorf("Error create response: %s", err.Error())
				continue
			}
			handler.JsonData = JsonData
			logrus.Warnf("Success update (dowload new data without support old)")

		case <-time.After(maxWaitPeriod):
			logrus.Warnf("This message appears every 10 seconds. General update every 30 seconds.")
		}
	}
}

func start() error {
	countriesMap, err := dataset.GetCountryMapAlpha2()
	if err != nil {
		logrus.Infof("Error create map with countries: %s", err.Error())
		return err
	}
	if err := InitDir(); err != nil {
		logrus.Infof("Error init path: %s", err.Error())
		return err
	}
	statusStruct := aggregator.StatusHandler{
		Status: true,
		Error:  "",
	}

	dataSms, dataMms, dataVoiceCall, dataEmail, dataBilling, dataSupport, dataIncident := getData()
	srv := new(aggregator.Server)
	var wg sync.WaitGroup

	wg.Add(6)
	startHandlers(&statusStruct, &wg, countriesMap, &dataSms, &dataMms, &dataVoiceCall, &dataEmail, &dataBilling, &dataSupport, &dataIncident)
	wg.Wait()

	JsonData, err := CreateResponse(&statusStruct, &dataSms, &dataMms, &dataVoiceCall, &dataEmail, &dataBilling, &dataSupport, &dataIncident)
	if err != nil {
		logrus.Infof("Error create response: %s", err.Error())
		return err
	}
	handler := web.Handler{
		JsonData: JsonData,
	}

	go func() {
		if err := srv.Run(port, handler.InitRoutes()); err != nil {
			logrus.Infof("Error occured while running http server")
		}
	}()

	go updateDataPeriodically(&statusStruct, &handler, countriesMap, &dataSms, &dataMms, &dataVoiceCall, &dataEmail, &dataBilling, &dataSupport, &dataIncident)

	logrus.Infof("Listening...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Infof("error occured on server shutting down: %s", err.Error())
		return err
	}
	return nil
}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := start(); err != nil {
		logrus.Fatalf("error in start: %s", err.Error())
		return
	}
}
