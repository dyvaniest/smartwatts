import React, { useState, useEffect } from "react";
import { List, Card, Spin, Tag, Layout, DatePicker, Pagination, Row, Col } from "antd";

function History() {
  const [data, setData] = useState([]);
  const [filteredData, setFilteredData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const token = localStorage.getItem("token")
        const response = await fetch("http://localhost:8080/data", {
          method: "GET",
          headers: {
              "Content-Type": "application/json",
              "Authorization":  `Bearer ${token}`,
          },
        });
        if (!response.ok) {
          throw new Error("Failed to fetch data");
        }
        const rawData = await response.json();

        // Filter dan urutkan data
        const processedData = processStatusChanges(rawData);

        setData(processedData);
        setFilteredData(processedData);
      } catch (error) {
        console.error("Error fetching data:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  // Fungsi untuk memproses perubahan status dan mengurutkan data
  const processStatusChanges = (rawData) => {
    // Mengelompokkan data berdasarkan Appliance
    const groupedData = rawData.reduce((acc, item) => {
      const { Appliance } = item;
      if (!acc[Appliance]) {
        acc[Appliance] = [];
      }
      acc[Appliance].push(item);
      return acc;
    }, {});

    const result = [];
    for (const appliance in groupedData) {
      const applianceData = groupedData[appliance];
      let previousStatus = null;

      // Filter hanya perubahan status
      applianceData.forEach((record) => {
        if (record.Status !== previousStatus) {
          result.push(record);
          previousStatus = record.Status;
        }
      });
    }

    // Mengurutkan hasil berdasarkan tanggal dan waktu
    return result.sort((a, b) => {
      const dateA = new Date(`${a.Date}T${a.Time}`);
      const dateB = new Date(`${b.Date}T${b.Time}`);
      return dateA - dateB;
    });
  };

  const handleDateFilter = (date, dateString) => {
    if (!dateString) {
      setFilteredData(data);
    } else {
      const filtered = data.filter(
        (item) => item.Date === dateString
      );
      setFilteredData(filtered);
      setCurrentPage(1);
    }
  };

  const handlePageChange = (page, pageSize) => {
    setCurrentPage(page);
    setPageSize(pageSize);
  };

  const paginatedData = filteredData.slice(
    (currentPage - 1) * pageSize,
    currentPage * pageSize
  );

  return (
    <Layout.Content style={{ padding: "10px", margin: "0 16px", marginBottom: 10 }}>
      <Card title="Appliances History" style={{ height: "100%", overflowY: "hidden" }}>
        {loading ? (
          <Spin size="large" style={{ display: "block", margin: "auto" }} />
        ) : (
          <>
            <Row style={{ marginBottom: 16 }}>
              <Col span={12}>
                <DatePicker onChange={handleDateFilter} style={{ width: "100%" }} />
              </Col>
            </Row>
            <List
              style={{ marginTop: 0 }}
              dataSource={paginatedData}
              renderItem={(item) => (
                <List.Item>
                  <List.Item.Meta
                    avatar={
                      <Tag
                        color={item.Status === "On" ? "green" : "red"}
                        style={{
                          width: 20,
                          height: 20,
                          borderRadius: "50%",
                          marginTop: 6,
                        }}
                      />
                    }
                    title={
                      <>
                        <span>{item.Appliance}</span>
                        <span style={{ fontSize: "12px", marginLeft: 8 }}>
                          • {item.Room}
                        </span>
                      </>
                    }
                    description={
                      <div style={{ display: "flex", justifyContent: "space-between" }}>
                        <span>{item.Status === "On" ? "Turn On" : "Turn Off"}</span>
                        <span>
                          {item.Date} {item.Time}
                        </span>
                      </div>
                    }
                  />
                </List.Item>
              )}
            />
            <Pagination
              style={{ marginTop: 16, textAlign: "center" }}
              current={currentPage}
              pageSize={pageSize}
              total={filteredData.length}
              onChange={handlePageChange}
              showSizeChanger
            />
          </>
        )}
      </Card>
    </Layout.Content>
  );
}

export default History;