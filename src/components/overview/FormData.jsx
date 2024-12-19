import React, { useState } from "react";
import { Layout, Form, Input, Button, Typography, message, DatePicker, TimePicker, InputNumber } from "antd";
import { useNavigate } from "react-router-dom";

const { Content } = Layout;
const { Title } = Typography;

function FormData() {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  // Fungsi untuk submit data
  const handleSubmit = async (values) => {
    setLoading(true);
    try {
      const token = localStorage.getItem("token"); // Mengambil token dari localStorage
      const payload = {
        date: values.date.format("YYYY-MM-DD"),
        time: values.time.format("HH:mm"),
        appliance: values.appliance,
        energy_consumption: values.energy_consumption,
        room: values.room,
        status: values.status,
      };

      const response = await fetch("http://localhost:8080/data", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(payload),
      });

      if (response.ok) {
        message.success("Data berhasil ditambahkan!");
        navigate("/overview"); // Navigasi ke halaman utama atau halaman lain setelah sukses
      } else {
        const errorData = await response.json();
        message.error(`Gagal menambahkan data: ${errorData.message || "Unknown error"}`);
      }
    } catch (error) {
      console.error("Error saat menambahkan data:", error);
      message.error("Terjadi kesalahan saat menambahkan data.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout style={{ minHeight: "100vh",  background: "none" }}>
      <Content
        style={{ display: "flex", justifyContent: "center", alignItems: "center" }}
      >
        <div
          style={{
            maxWidth: 500,
            width: "100%",
            padding: 24,
            background: "#fff",
            borderRadius: 8,
            boxShadow: "0 2px 8px rgba(0, 0, 0, 0.1)",
          }}
        >
          <Title level={3} style={{ textAlign: "center", fontFamily: "Rubik" }}>
            Add Appliance Data
          </Title>
          <Form
            layout="vertical"
            onFinish={handleSubmit}
            style={{ marginTop: 16 }}
          >
            {/* Input Tanggal */}
            <Form.Item
              label="Date"
              name="date"
              rules={[{ required: true, message: "Please select the date!" }]}
            >
              <DatePicker style={{ width: "100%" }} />
            </Form.Item>

            {/* Input Waktu */}
            <Form.Item
              label="Time"
              name="time"
              rules={[{ required: true, message: "Please select the time!" }]}
            >
              <TimePicker style={{ width: "100%" }} format="HH:mm" />
            </Form.Item>

            {/* Input Nama Appliance */}
            <Form.Item
              label="Appliance Name"
              name="appliance"
              rules={[{ required: true, message: "Please input the appliance name!" }]}
            >
              <Input placeholder="Enter appliance name" />
            </Form.Item>

            {/* Input Energy Consumption */}
            <Form.Item
              label="Energy Consumption (kWh)"
              name="energy_consumption"
              rules={[
                { required: true, message: "Please input the energy consumption!" },
                { type: "number", min: 0, message: "Value must be positive!" },
              ]}
            >
              <InputNumber placeholder="Enter energy consumption" style={{ width: "100%" }} />
            </Form.Item>

            {/* Input Room */}
            <Form.Item
              label="Room"
              name="room"
              rules={[{ required: true, message: "Please input the room!" }]}
            >
              <Input placeholder="Enter room name" />
            </Form.Item>

            {/* Input Status */}
            <Form.Item
              label="Status"
              name="status"
              rules={[
                {
                  required: true,
                  message: "Please input the status (On/Off)!",
                },
                {
                  validator: (_, value) =>
                    ["On", "Off"].includes(value)
                      ? Promise.resolve()
                      : Promise.reject(new Error("Status must be 'On' or 'Off'!")),
                },
              ]}
            >
              <Input placeholder="Enter status (e.g., On or Off)" />
            </Form.Item>

            {/* Tombol Submit */}
            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                block
                loading={loading}
              >
                Submit
              </Button>
            </Form.Item>
          </Form>
        </div>
      </Content>
    </Layout>
  );
}

export default FormData;
