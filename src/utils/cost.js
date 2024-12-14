import React from "react";
import { Card, Typography } from "antd";
import { DollarCircleOutlined } from "@ant-design/icons";

const { Title } = Typography;

const EstimatedCostCard = ({ totalConsumption, estimatedCost }) => {
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

export default EstimatedCostCard;
