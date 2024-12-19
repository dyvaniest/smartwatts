import React, { useState } from "react";
import { Layout, Form, Input, Button, Card, Spin, List, message } from "antd";
import { SendOutlined } from "@ant-design/icons";

import { formatResponse } from "../utils/cost";

const Chat = () => {
  const [chatQuery, setChatQuery] = useState("");
  const [chatResponses, setChatResponses] = useState([]);
  const [chatLoading, setChatLoading] = useState(false);

  // Handle Chat Submission
  const handleChatSubmit = async () => {
    if (!chatQuery) {
      message.error("Please enter a message to chat with AI.");
      return;
    }
    setChatLoading(true);
    try {
      const token = localStorage.getItem("token")
      const response = await fetch("http://localhost:8080/chat", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({ query: chatQuery }),
      });

      if (!response.ok) {
        throw new Error("Failed to chat with AI");
      }

      const data = await response.json();
      setChatResponses((prev) => [
        ...prev,
        { user: chatQuery, ai: data.answer },
      ]);
      setChatQuery("");
    } catch (error) {
      message.error("Error: " + error.message);
    } finally {
      setChatLoading(false);
    }
  };

  return (
    <Layout.Content style={{ padding: "10px", margin: "0 16px" }}>
      <Card title="Chat with AI" style={{ height: "100%", overflowY: "hidden" }}>
        {/* Chat History */}
        <div
          style={{
            maxHeight: "400px",
            overflowY: "auto",
            marginBottom: "16px",
            borderRadius: "8px",
          }}
        >
          {chatResponses.map((item, index) => (
            <div key={index} style={{ marginBottom: "16px" }}>
              {/* User Message */}
              <div
                style={{
                  textAlign: "right",
                  marginBottom: "8px",
                }}
              >
                <div
                  style={{
                    display: "inline-block",
                    background: "#1890ff",
                    color: "#fff",
                    padding: "8px 12px",
                    borderRadius: "16px",
                    maxWidth: "70%",
                    wordWrap: "break-word",
                  }}
                >
                  {item.user}
                </div>
              </div>
              {/* AI Response */}
              <div
                style={{
                  textAlign: "left",
                }}
              >
                <div
                  style={{
                    display: "inline-block",
                    background: "#f0f0f0",
                    color: "#000",
                    padding: "8px 12px",
                    borderRadius: "16px",
                    maxWidth: "70%",
                    wordWrap: "break-word",
                  }}
                >
                  {formatResponse(item.ai)}
                </div>
              </div>
            </div>
          ))}
          {chatLoading && (
            <div style={{ textAlign: "center" }}>
              <Spin />
            </div>
          )}
        </div>
        {/* Chat Input */}
        <Form 
          layout="inline" 
          onFinish={handleChatSubmit} 
          style={{ position: "relative", bottom: 0, left: 0, right:0 }}
        >
          <Form.Item style={{ flexGrow: 1, marginRight: "8px" }}>
            <Input
              placeholder="Type your message..."
              value={chatQuery}
              onChange={(e) => setChatQuery(e.target.value)}
              disabled={chatLoading}
              style={{
                borderRadius: 20,
                height: "50px"
              }}
            />
          </Form.Item>
          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={chatLoading}
              disabled={chatLoading}
              style={{ height: "45px", borderRadius: 60}}
            >
              <SendOutlined/>
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </Layout.Content>
  );
};

export default Chat;
