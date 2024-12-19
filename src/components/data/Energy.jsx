import React, { useState, useEffect } from "react";
import { Row, Col, Typography, Select, Card, Progress } from "antd";
import { Line } from "@ant-design/plots";
import { EstimatedCostCard } from "../../utils/cost";

const { Title } = Typography;
const { Option } = Select;

function Energy() {
  const [energyData, setEnergyData] = useState(null);
  const [additionalData, setAdditionalData] = useState([]);
  const [timeRange, setTimeRange] = useState("day"); // Default time range

  useEffect(() => {
    const fetchData = async () => {
      try {
        const token = localStorage.getItem("token");
        const response = await fetch("http://localhost:8080/analytics-energy", {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
        });
        const data = await response.json();
        setEnergyData(data);
      } catch (error) {
        console.error("Error fetching energy data:", error);
      }
    };

    const fetchAdditionalData = async () => {
      try {
        const token = localStorage.getItem("token");
        const response = await fetch("http://localhost:8080/data", {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
        });
        const data = await response.json();

        // Validate and normalize additional data
        const normalizedData = data.map((item) => ({
          time: `${item.date} ${item.time}`, // Combine date and time
          value: item.energy_consumption, // Energy consumption
          room: item.room, // Room name
        }));

        setAdditionalData(normalizedData);
      } catch (error) {
        console.error("Error fetching additional data:", error);
      }
    };

    fetchData();
    fetchAdditionalData();
  }, []);

  const handleTimeRangeChange = (value) => {
    setTimeRange(value);
  };

  if (!energyData) {
    return <div>Loading...</div>;
  }

  const chartData = [];
  const colors = ["#1890ff", "#52c41a", "#faad14", "#ff4d4f", "#722ed1"];

  // Process energy data based on selected time range
  Object.entries(
    timeRange === "day"
      ? energyData.daily_consumption
      : timeRange === "week"
      ? energyData.weekly_consumption
      : timeRange === "month"
      ? energyData.monthly_consumption
      : {}
  ).forEach(([room, consumption], index) => {
    if (Array.isArray(consumption)) {
      consumption.forEach((dataPoint) => {
        chartData.push({
          time: dataPoint.time,
          value: dataPoint.value,
          room: room,
        });
      });
    }
  });

  // Add additional data to chart
  additionalData.forEach((item) => {
    chartData.push({
      time: item.time,
      value: item.value,
      room: item.room,
    });
  });

  // Calculate total energy consumption
  const totalConsumption = chartData.reduce((acc, item) => acc + item.value, 0);

  const chartConfig = {
    data: chartData,
    xField: "time",
    yField: "value",
    seriesField: "room",
    color: ({ room }) =>
      colors[Object.keys(energyData.room_consumption).indexOf(room) % colors.length],
    smooth: true,
    height: 300,
    xAxis: {
      title: { text: "Time" },
    },
    yAxis: {
      title: { text: "Energy (kWh)" },
    },
    legend: {
      position: "top",
    },
  };

  return (
    <div style={{ padding: "20px", backgroundColor: "white", borderRadius: 10 }}>
      <Title level={3} style={{ fontSize: "14pt", fontFamily: "Rubik" }}>
        Energy Consumption
      </Title>
      <Row gutter={[32, 16]}>
        <Col span={8}>
          <EstimatedCostCard
            totalConsumption={totalConsumption}
            estimatedCost={energyData.estimated_cost}
          />
        </Col>
      </Row>

      {/* Room Energy Consumption */}
      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        <Col span={8}>
          <Row gutter={[16, 16]}>
            {Object.entries(energyData.room_consumption).map(([room, consumption], index) => (
              <Col span={12} key={room}>
                <Card style={{ padding: 0 }}>
                  <Row gutter={[16, 16]} style={{ margin: 2 }}>
                    <Col span={4}>
                      <Progress
                        type="circle"
                        percent={((consumption / totalConsumption) * 100).toFixed(2)} // Normalization for progress
                        width={50}
                        strokeColor={["#52c41a", "#ff4d4f", "#722ed1", "#faad14"][index % 4]}
                      />
                    </Col>
                    <Col span={24} style={{ paddingLeft: 10 }}>
                      <Title level={5} style={{ fontSize: "12pt", fontFamily: "Rubik" }}>
                        {room}
                      </Title>
                      <p style={{ marginTop: 0 }}>{consumption.toFixed(2)} kWh</p>
                    </Col>
                  </Row>
                </Card>
              </Col>
            ))}
          </Row>
        </Col>

        {/* Chart */}
        <Col span={16}>
          <Card>
            <Row justify="space-between" align="middle">
              <Col>
                <Title level={4} style={{ fontSize: "12pt", fontFamily: "Rubik" }}>
                  Energy Consumption Chart
                </Title>
              </Col>
              <Select defaultValue="day" style={{ width: 120 }} onChange={handleTimeRangeChange}>
                <Option value="day">Day</Option>
                <Option value="week">Week</Option>
                <Option value="month">Month</Option>
              </Select>
            </Row>
            <Line {...chartConfig} />
          </Card>
        </Col>
      </Row>
      <Row>
          <div>{energyData.ai_suggestion}</div>
      </Row>
    </div>
  );
}

export default Energy;
