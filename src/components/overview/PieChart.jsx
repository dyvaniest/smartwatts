import React, { useState, useEffect } from "react";
import { Pie } from "@ant-design/plots";
import { Card, Spin, message } from "antd";
import { Link } from "react-router-dom";

function PieChart() {
    const [energyData, setEnergyData] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
        try {
            const token = localStorage.getItem("token");
            const response = await fetch("http://localhost:8080/analytics-energy", {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization":  `Bearer ${token}`,
                },

            });

            if (!response.ok) {
            throw new Error("Failed to fetch data");
            }
            const data = await response.json();
            setEnergyData(data);
            setLoading(false);
        } catch (error) {
            console.error("Error fetching data:", error);
            message.error("Failed to load data. Please try again later.");
            setLoading(false);
        }
        };
        fetchData();
    }, []);

    if (loading) {
        return (
        <div style={{ display: "flex", justifyContent: "center", alignItems: "center", height: "100vh" }}>
            <Spin size="large" />
        </div>
        );
    }

    if (!energyData) {
        return <p>No data available</p>;
    }

    const pieChartData = Object.entries(energyData.room_consumption).map(
        ([room, consumption]) => ({
          type: room,
          value: consumption,
        })
      );
    
    const pieConfig = {
        appendPadding: 8,
        data: pieChartData,
        angleField: "value",
        colorField: "type",
        radius: 1,
        innerRadius: 0.7,
        label: {
            type: "outer",
            content: ({ percent }) => `${(percent * 100).toFixed(2)}%`,
            style: {
                textAlign: "center",
                fontSize: 10,
            },
        },
        interactions: [
            {
                type: "element-active",
            },
        ],
        legend: {
            position: "right",
        },
        color: [
            "#ff4d4f",
            "#ffa940",
            "#ffc53d",
            "#73d13d",
            "#40a9ff",
            "#9254de",
        ],
        width: 200,
        height: 200,
    };
      
    return (
        <div>
        <Card title="Consumption by Room" extra={<Link to="/analytics">View analytics...</Link>} style={{ border: "none" }}>
            <Pie {...pieConfig} />
        </Card>
        </div>
    );
}
      
export default PieChart;
      