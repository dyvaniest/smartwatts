import React, { useState } from "react";
import { Layout, Breadcrumb, Avatar, Row, Col, Modal, Dropdown, message, Menu } from "antd";
import "../index.css";
import { BellOutlined, HistoryOutlined, SettingOutlined, UserOutlined } from "@ant-design/icons";
import History from "../pages/History";

function Header() {
  const [isModalVisible, setIsModalVisible] = useState(false); // Kontrol visibilitas modal

  const handleMenuClick = async ({ key }) => {
    if (key === "logout") {
      try {
        const token = localStorage.getItem("token");
        if (!token) {
          throw new Error("No auth token found");
        }

        const response = await fetch("http://localhost:8080/logout", {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
          },
          credentials: "include",
        });

        if (!response.ok) {
          const error = await response.text();
          console.error("Backend error:", error);
          throw new Error("Logout failed");
        }

        message.success("Logged out successfully!");
        localStorage.removeItem("token");
        window.location.href = "/";
      } catch (error) {
        message.error("Error logging out. Please try again.");
        console.error("Logout error:", error);
      }
    }
  };

  const handleOpenModal = () => {
    setIsModalVisible(true);
  };

  const handleCloseModal = () => {
    setIsModalVisible(false);
  };

  const menu = (
    <Menu onClick={handleMenuClick}>
      <Menu.Item key="logout">Logout</Menu.Item>
    </Menu>
  );

  return (
    <Layout.Header
      style={{
        background: "none",
        padding: "0 24px",
      }}
    >
      <Row align="middle" justify="space-between" gutter={[8, 8]}>
        <Col>
          <Breadcrumb>
            <Breadcrumb.Item>Home</Breadcrumb.Item>
            <Breadcrumb.Item>Dashboard</Breadcrumb.Item>
          </Breadcrumb>
        </Col>
        <Col>
          <Row align="middle" gutter={16}>
            <Col>
              <SettingOutlined style={{ fontSize: "20px", color: "#112d3c" }} />
            </Col>
            <Col>
              <HistoryOutlined
                style={{
                  fontSize: "20px",
                  color: "#112d3c",
                  cursor: "pointer",
                }}
                onClick={handleOpenModal}
              />
            </Col>
            <Col>
              <Dropdown overlay={menu} placement="bottomRight">
                <Avatar
                  size={32}
                  icon={<UserOutlined />}
                  style={{ backgroundColor: "#112d3c", cursor: "pointer" }}
                />
              </Dropdown>
            </Col>
          </Row>
        </Col>
      </Row>

      <Modal
        visible={isModalVisible}
        onCancel={handleCloseModal} 
        footer={null}
        width={1000}
      >
        <History /> 
      </Modal>
    </Layout.Header>
  );
}

export default Header;
