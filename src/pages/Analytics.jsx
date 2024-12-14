import React from "react";
import { Layout, Row, Col, message } from "antd";
import Energy from "../components/data/Energy";
import AnalyzeAi from "../components/data/Analyze";

const Analytics = () => {
  return (
    <Layout.Content style={{ padding: "10px", margin: "0 16px" }}>
      <Row gutter={[16, 16]} style={{ marginTop: 0, marginBottom: 20 }}>
        <Col span={24}>
          <AnalyzeAi/>
        </Col>
        <Col span={24}>
          <Energy/>
        </Col>
      </Row>
      
    </Layout.Content>
  );
};

export default Analytics;
