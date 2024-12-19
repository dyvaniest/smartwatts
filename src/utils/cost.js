import React from "react";
import { Card, Typography } from "antd";
import { DollarCircleOutlined } from "@ant-design/icons";

const { Title } = Typography;

export const EstimatedCostCard = ({ totalConsumption, estimatedCost }) => {
  return (
    <Card
      title="Estimated Cost"
      bordered={false}
      extra={<DollarCircleOutlined style={{ fontSize: 24 }} />}
    >
      <p>Total Energy: {totalConsumption.toFixed(2)} kWh </p>
      <Title level={4} style={{ fontFamily: "Rubik", fontSize: "18px" }}>
        {estimatedCost.toLocaleString()} IDR
      </Title>
    </Card>
  );
};

// export default EstimatedCostCard;

// Format AI response for better display
export const formatResponse = (response) => {
  const parts = response.split(/(\*\*[^*]+\*\*|\n)/g); // Split by **bold** or \n
  return parts.map((part, index) => {
    if (part.startsWith("**") && part.endsWith("**")) {
      return (
        <strong key={index}>{part.slice(2, -2)}</strong> // Render bold text
      );
    } else if (part === "\n") {
      return <br key={index} />; // Render line break
    } else {
      return part; // Render plain text
    }
  });
};