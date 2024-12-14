import React, { useState, useEffect } from "react";
import { Card, Spin, message } from "antd";
import { Line } from "@ant-design/plots";

function Realtime() {
  const [energyData, setEnergyData] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const token = localStorage.getItem("token")
        const energyResponse = await fetch("http://localhost:8080/data", {
          method: "GET",
          headers: {
              "Content-Type": "application/json",
              "Authorization":  `Bearer ${token}`,
          },
        });
      
        const rawEnergyData = await energyResponse.json();

        const groupedEnergyData = rawEnergyData.reduce((acc, curr) => {
          const { Time, EnergyConsumption } = curr;
          acc[Time] = (acc[Time] || 0) + EnergyConsumption;
          return acc;
        }, {});

        // Convert to chart-ready format
        const chartData = Object.keys(groupedEnergyData).map((time) => ({
          time,
          value: groupedEnergyData[time],
        }));

        setEnergyData(chartData);
        setLoading(false);
      } catch (error) {
        console.error("Error fetching data:", error);
        message.error("Failed to load data. Please try again later.");
        setLoading(false);
      }
    };

    fetchData();

    // Refresh data every 5 seconds
    const interval = setInterval(fetchData, 5000);
    return () => clearInterval(interval);
  }, []);

  const lineConfig = {
    data: energyData,
    xField: "time",
    yField: "value",
    smooth: true,
    height: 200,
    xAxis: {
      title: { text: "Time" },
      label: {
        formatter: (val) => `${val}`,
      },
    },
    yAxis: {
      title: { text: "Energy Consumption (kWh)" },
    },
    color: "#1890ff",
    tooltip: {
      showMarkers: true,
    },
    animation: {
      appear: {
        animation: 'path-in',
        duration: 1000,
      },
    },
  };

  return (
    <div>
      {loading ? (
        <div style={{ display: "flex", justifyContent: "center", alignItems: "center", height: "100vh" }}>
          <Spin size="large" />
        </div>
      ) : (
        <Card title="Real-Time Energy Consumption" style={{ border: "none" }}>
          <Line {...lineConfig} />
        </Card>
      )}
    </div>
  );
}

export default Realtime;
