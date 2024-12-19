package service

import (
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

type FileService struct {
	Repo *repository.FileRepository
}

type EnergyData struct {
	Date              string  `json:"date"`
	Time              string  `json:"time"`
	Appliance         string  `json:"appliance"`
	EnergyConsumption float64 `json:"energy_consumption"`
	Room              string  `json:"room"`
	Status            string  `json:"status"`
}

// ReadCSV reads data from a CSV file and converts it into a slice of EnergyData.
func (f *FileService) ReadCSV(filePath string) ([]EnergyData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []EnergyData
	for i, row := range rows {
		if i == 0 {
			continue
		}

		energyConsumption, _ := strconv.ParseFloat(row[3], 64)
		data = append(data, EnergyData{
			Date:              row[0],
			Time:              row[1],
			Appliance:         row[2],
			EnergyConsumption: energyConsumption,
			Room:              row[4],
			Status:            row[5],
		})
	}

	return data, nil
}

// WriteCSV writes a slice of EnergyData back to a CSV file.
func (f *FileService) WriteCSV(filePath string, data []EnergyData) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"Date", "Time", "Appliance", "Energy_Consumption", "Room", "Status"})

	// Write rows
	for _, record := range data {
		writer.Write([]string{
			record.Date,
			record.Time,
			record.Appliance,
			strconv.FormatFloat(record.EnergyConsumption, 'f', 2, 64),
			record.Room,
			record.Status,
		})
	}

	return nil
}

func (f *FileService) AddDataCSV(filePath string, data EnergyData) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	w := csv.NewWriter(file)
	defer w.Flush()

	// Write rows
	w.Write([]string{
		data.Date,
		data.Time,
		data.Appliance,
		strconv.FormatFloat(data.EnergyConsumption, 'f', 2, 64),
		data.Room,
		data.Status,
	})

	return nil
}

// AnalyzeEnergyConsumption calculates the average energy consumption per day, week, and month.
func (f *FileService) AnalyzeEnergyConsumption(data []EnergyData) (map[string]float64, map[string]float64, map[string]float64) {
	dailyConsumption := make(map[string]float64)
	weeklyConsumption := make(map[string]float64)
	monthlyConsumption := make(map[string]float64)

	for _, record := range data {
		date, _ := time.Parse("2006-01-02", record.Date)
		year, month, day := date.Date()
		_, week := date.ISOWeek()

		dailyKey := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
		weeklyKey := fmt.Sprintf("%04d-W%02d", year, week)
		monthlyKey := fmt.Sprintf("%04d-%02d", year, month)

		dailyConsumption[dailyKey] += record.EnergyConsumption
		weeklyConsumption[weeklyKey] += record.EnergyConsumption
		monthlyConsumption[monthlyKey] += record.EnergyConsumption
	}

	for key := range dailyConsumption {
		dailyConsumption[key] = round(dailyConsumption[key]/24, 2)
	}

	for key := range weeklyConsumption {
		weeklyConsumption[key] = round(weeklyConsumption[key]/168, 2)
	}

	for key := range monthlyConsumption {
		monthlyConsumption[key] = round(monthlyConsumption[key]/720, 2)
	}

	return dailyConsumption, weeklyConsumption, monthlyConsumption
}

// CalculateRoomEnergyConsumption calculates the total energy consumption for each room.
func (f *FileService) CalculateRoomEnergyConsumption(data []EnergyData) map[string]float64 {
	roomConsumption := make(map[string]float64)

	for _, record := range data {
		roomConsumption[record.Room] += record.EnergyConsumption
	}

	for key := range roomConsumption {
		roomConsumption[key] = round(roomConsumption[key], 2)
	}

	return roomConsumption
}

// EstimateElectricityCost estimates the total electricity cost based on the energy consumption.
func (f *FileService) EstimateElectricityCost(data []EnergyData, costPerKWh float64) float64 {
	totalConsumption := 0.0

	for _, record := range data {
		totalConsumption += record.EnergyConsumption
	}

	return round(totalConsumption*costPerKWh, 2)
}

// round membulatkan nilai ke dua angka di belakang koma
func round(val float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Round(val*p) / p
}
