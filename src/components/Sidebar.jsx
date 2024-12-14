import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import Icon, {
  AppstoreOutlined,
  AreaChartOutlined,
  FileSearchOutlined,
  HistoryOutlined,
  MessageOutlined,
} from "@ant-design/icons";
import { Layout, Menu, Image } from "antd";
import Logo from "../assets/smartwatts-logo.png";
import Slogan from "../assets/smartWatts.png";

const getItem = (label, key, icon, path) => ({
  key,
  icon,
  label,
  path,
});

const items = [
  getItem("Overview", "1", <AppstoreOutlined />, "/overview"),
  getItem("Analytics", "2", <AreaChartOutlined />, "/analytics"),
  getItem("Chat", "3", <MessageOutlined />, "/chat"),
  getItem("Recommendations", "4", <FileSearchOutlined/>, "/recommendations"),
];

const Sidebar = ({ collapsed, onCollapse }) => {
  const navigate = useNavigate();

  const handleMenuClick = ({ key }) => {
    const item = items.find((i) => i.key === key);
    if (item && item.path) {
      navigate(item.path);
    }
  };

  return (
    <Layout.Sider
      style={{ height: "100vh", position: "fixed", left: 0, top: 0 }}
      collapsible
      collapsed={collapsed}
      onCollapse={onCollapse}
      theme="light"
    >
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: collapsed ? "center" : "flex-start",
          padding: "8px",
        }}
      >
        <Image
          src={Logo}
          style={{
            width: 50,
            height: 50,
          }}
          preview={false}
        />
        {!collapsed && (
          <Image
            src={Slogan}
            style={{
              marginLeft: 8,
              width: 120,
            }}
            preview={false}
          />
        )}
      </div>
      <Menu
        defaultSelectedKeys={["1"]}
        mode="inline"
        items={items}
        theme="light"
        onClick={handleMenuClick}
      />
    </Layout.Sider>
  );
};

export default Sidebar;
