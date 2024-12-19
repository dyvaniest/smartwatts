import React, { useState, useEffect } from "react";
import { Layout, Typography, Row, Col, Card, Modal, FloatButton } from "antd";
import Rooms from "../components/Rooms";
import "../index.css";
import Realtime from "../components/overview/Realtime";
import PieChart from "../components/overview/PieChart";
import EstimatedCostCard from "../utils/cost";
import { PlusOutlined } from "@ant-design/icons";
import FormData from "../components/overview/FormData";

function Overview() {
  const [energyData, setEnergyData] = useState(null);
  const [isModalVisible, setIsModalVisible] = useState(false); // Kontrol popup

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

    fetchData();
  }, []);

  if (!energyData) {
    return <div>Loading...</div>;
  }

  const totalConsumption = Object.values(energyData.room_consumption).reduce(
    (acc, value) => acc + value,
    0
  );

  return (
    <Layout.Content style={{ margin: "0 16px" }}>
      <div
        style={{
          padding: 10,
          minHeight: 360,
        }}
      >
        {/* Float Button untuk membuka Chat */}
        <FloatButton
          icon={<PlusOutlined />}
          type="primary"
          tooltip={<div>Add Appliance</div>}
          style={{
            insetInlineEnd: 24,
            width: 50,
            height: 50,
          }}
          onClick={() => setIsModalVisible(true)} // Tampilkan popup saat diklik
        />

        {/* Popup Modal untuk Chat */}
        <Modal
          open={isModalVisible}
          onCancel={() => setIsModalVisible(false)}
          footer={null}
          style={{ width: '100vh'}}
          
        >
          <div style={{ padding: 0}}>
              <FormData />
          </div>
        </Modal>

        <Rooms />
        <Typography.Title
          style={{ fontSize: "18px", marginTop: 12, fontFamily: "Rubik" }}
        >
          Usage
        </Typography.Title>
        <Card>
          <Row gutter={[32, 16]} style={{ marginTop: 16 }}>
            <Col span={8}>
              <EstimatedCostCard
                totalConsumption={totalConsumption}
                estimatedCost={energyData.estimated_cost}
              />
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Realtime />
            </Col>
            <Col span={12}>
              <PieChart />
            </Col>
          </Row>
        </Card>
      </div>
    </Layout.Content>
  );
}

export default Overview;
