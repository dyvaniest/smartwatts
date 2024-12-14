import React, { useState } from "react";
import { Layout, Form, Input, Button, Card, Spin, message, Select } from "antd";

const { Option } = Select;

const Analyze = () => {
  const [selectedOption, setSelectedOption] = useState("");
  const [customQuestion, setCustomQuestion] = useState("");
  const [analyzeResult, setAnalyzeResult] = useState("");
  const [analyzeLoading, setAnalyzeLoading] = useState(false);

  const presetQuestions = [
    "What the appliance consumes the most energy?",
    "Which the appliance consumes the least energy?",
    "How much energy is consumed by the most wasteful appliance?"
  ];

  const handleAnalyzeSubmit = async () => {
    const analyzeQuery =
      selectedOption === "Other" ? customQuestion : selectedOption;

    if (!analyzeQuery) {
      message.error("Please select or enter a query for analysis.");
      return;
    }

    setAnalyzeLoading(true);
    try {
      const token = localStorage.getItem("token");
      const response = await fetch("http://localhost:8080/analyze-ai", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ query: analyzeQuery }), 
      });

      if (!response.ok) {
        throw new Error("Failed to analyze data.");
      }

      const data = await response.json();
      setAnalyzeResult(data.result);
      message.success("Analysis completed!");
    } catch (error) {
      message.error("Error: " + error.message);
    } finally {
      setAnalyzeLoading(false);
    }
  };

  const handleOptionChange = (value) => {
    setSelectedOption(value);
    if (value !== "Other") {
      setCustomQuestion("");
    }
  };

  return (
    <>
      {/* Analyze Section */}
      <Card title="Analyze Data" style={{ marginBottom: "20px" }}>
        <Form layout="vertical" onFinish={handleAnalyzeSubmit}>
          {/* Dropdown untuk preset pertanyaan */}
          <Form.Item label="Select a Question:">
            <Select
              placeholder="Select a question"
              onChange={handleOptionChange}
              value={selectedOption || undefined}
              allowClear
            >
              {presetQuestions.map((question, index) => (
                <Option key={index} value={question}>
                  {question}
                </Option>
              ))}
              <Option value="Other">Other question</Option>
            </Select>
          </Form.Item>

          {/* Input manual muncul hanya jika Other dipilih */}
          {selectedOption === "Other" && (
            <Form.Item label="Enter your question:">
              <Input
                placeholder="Enter your custom question"
                value={customQuestion}
                onChange={(e) => setCustomQuestion(e.target.value)}
              />
            </Form.Item>
          )}

          {/* Tombol Analyze */}
          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={analyzeLoading}
              style={{ width: "100%" }}
            >
              Analyze
            </Button>
          </Form.Item>
        </Form>

        {/* Hasil Analisis */}
        {analyzeLoading ? (
          <Spin />
        ) : (
          analyzeResult && (
            <Card type="inner" title="Analysis Result" style={{ marginTop: 24 }}>
              {analyzeResult}
            </Card>
          )
        )}
      </Card>
    </>
  );
};

export default Analyze;
