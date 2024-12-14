import React, { useState } from "react";
import { Layout, Form, Input, Button, Typography, message } from "antd";
import { useNavigate } from "react-router-dom";

const { Content } = Layout;
const { Title } = Typography;

function Home({ isLogin, onToggle }) {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (values) => {
    setLoading(true);
    const endpoint = isLogin ? "http://localhost:8080/login" : "http://localhost:8080/register";
    try {
      const response = await fetch(endpoint, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(values),
      });

      if (!response.ok) {
        throw new Error("Request failed");
      }

      const result = await response.json();
      localStorage.setItem("token", result.token);
      message.success(isLogin ? "Login successful!" : "Registration successful!");

      if (isLogin) {
        navigate("/overview");
      }

      console.log(result);

    } catch (error) {
      message.error("Something went wrong. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout style={{ minHeight: "100vh", background: "#f5f5f5" }}>
      <Content style={{ display: "flex", justifyContent: "center", alignItems: "center" }}>
        <div
          style={{
            maxWidth: 400,
            width: "100%",
            padding: 24,
            background: "#fff",
            borderRadius: 8,
            boxShadow: "0 2px 8px rgba(0, 0, 0, 0.1)",
          }}
        >
          <Title level={3} style={{ textAlign: "center" }}>
            {isLogin ? "Sign In" : "Sign Up"}
          </Title>
          <Form
            layout="vertical"
            onFinish={handleSubmit}
            style={{ marginTop: 16 }}
          >
            <Form.Item
              name="email"
              label="Email"
              rules={[{ required: true, message: "Please enter your email" }]}
            >
              <Input type="email" placeholder="Enter your email" />
            </Form.Item>

            <Form.Item
              name="password"
              label="Password"
              rules={[{ required: true, message: "Please enter your password" }]}
            >
              <Input.Password placeholder="Enter your password" />
            </Form.Item>

            {!isLogin && (
              <>
                <Form.Item
                  name="confirmPassword"
                  label="Confirm Password"
                  dependencies={["password"]}
                  rules={[
                    { required: true, message: "Please confirm your password" },
                    ({ getFieldValue }) => ({
                      validator(_, value) {
                        if (!value || getFieldValue("password") === value) {
                          return Promise.resolve();
                        }
                        return Promise.reject(new Error("Passwords do not match!"));
                      },
                    }),
                  ]}
                >
                  <Input.Password placeholder="Confirm your password" />
                </Form.Item>

                <Form.Item
                  name="fullName"
                  label="Full Name"
                  rules={[{ required: true, message: "Please enter your full name" }]}
                >
                  <Input placeholder="Enter your full name" />
                </Form.Item>
              </>
            )}

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                loading={loading}
                block
              >
                {isLogin ? "Sign In" : "Sign Up"}
              </Button>
            </Form.Item>
          </Form>

          <div style={{ textAlign: "center", marginTop: 16 }}>
            {isLogin ? (
              <>
                Don’t have an account? <a onClick={() => onToggle(false)}>Sign up</a>
              </>
            ) : (
              <>
                Already have an account? <a onClick={() => onToggle(true)}>Login</a>
              </>
            )}
          </div>
        </div>
      </Content>
    </Layout>
  );
}

export default Home;
