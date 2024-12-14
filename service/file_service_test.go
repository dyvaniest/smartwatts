package service_test

import (
	"a21hc3NpZ25tZW50/service"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FileService", func() {
	var (
		fileService *service.FileService
		mockData    []service.EnergyData
		fileTest    string
	)

	BeforeEach(func() {
		fileService = &service.FileService{}
		fileTest = "test_data.csv"

		mockData = []service.EnergyData{
			{Date: "2024-12-10", Time: "11:00", Appliance: "Air Conditioner", EnergyConsumption: 2.45, Room: "Bedroom", Status: "On"},
			{Date: "2024-12-10", Time: "12:00", Appliance: "TV", EnergyConsumption: 0.8, Room: "Living Room", Status: "Off"},
		}

		err := fileService.WriteCSV(fileTest, mockData)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		err := os.Remove(fileTest)
		Expect(err).To(BeNil())
	})

	Describe("ReadCSV", func() {
		It("should read data from a CSV file correctly", func() {
			data, err := fileService.ReadCSV(fileTest)
			Expect(err).ToNot(HaveOccurred())
			Expect(data).To(Equal(mockData))
		})

	})

	Describe("WriteCSV", func() {
		It("should write data to a CSV file correctly", func() {
			newData := []service.EnergyData{
				{Date: "2024-12-10", Time: "13:00", Appliance: "EVCar", EnergyConsumption: 99.9, Room: "Garage", Status: "Off"},
			}

			err := fileService.WriteCSV(fileTest, newData)
			Expect(err).To(BeNil())

			data, err := fileService.ReadCSV(fileTest)
			Expect(err).To(BeNil())
			Expect(data).To(Equal(newData))
		})
	})
})
