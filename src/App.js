import React, { useState } from "react";
import { Routes, Route } from "react-router-dom";
import { Layout } from "antd";
import "./index.css";
import Sidebar from "./components/Sidebar";
import Home from "./pages/Home";
import Overview from "./pages/Overview";
import History from "./pages/History";
import Analytics from "./pages/Analytics";
import Chat from "./pages/Chat";
import Footer from "./components/Footer";
import Header from "./components/Header";
import Recommendations from "./pages/Recom";

const App = () => {
  const [isLogin, setIsLogin] = useState(true);
  const [collapsed, setCollapsed] = useState(false);

  return (
    <Layout style={{ minHeight: "100vh", backgroundColor: "#E8F3FC" }}>
      <Routes>
        {/* Route untuk halaman Home */}
        <Route path="/" element={<Home isLogin={isLogin} onToggle={setIsLogin}/>} />

        {/* Routes untuk halaman lainnya dengan sidebar dan header */}
        <Route
          path="/*"
          element={
            <Layout style={{ minHeight: "100vh", backgroundColor: "#E8F3FC" }}>
              <Sidebar 
                collapsed={collapsed}
                onCollapse={(value) => setCollapsed(value)}
              />
      
              <Layout style={{
                  backgroundColor: "#E8F3FC",
                  marginLeft: collapsed ? 80 : 200,
                  transition: "margin-left 0.2s ease-in-out",
                }}>
                <Header />
                <Routes>
                  <Route path="/overview" element={<Overview />} />
                  <Route path="/analytics" element={<Analytics />} />
                  <Route path="/chat" element={<Chat />} />
                  <Route path="/recommendations" element={<Recommendations />} />
                </Routes>
                <Footer />
              </Layout>
            </Layout>
          }
        />
      </Routes>
    </Layout>
  );
};

export default App;
