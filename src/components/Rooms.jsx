import React, { useState, useEffect } from "react";
import { Card, Row, Col, Typography, Spin, Switch, Tabs } from "antd";

const { TabPane } = Tabs;

function Rooms() {
  const [rooms, setRooms] = useState([]);
  const [appliances, setAppliances] = useState({});
  const [selectedRoom, setSelectedRoom] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const token = localStorage.getItem("token");
        const response = await fetch("http://localhost:8080/data", {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`,
          },
        });

        if (!response.ok) {
          throw new Error("Network response was not ok");
        }

        const rawData = await response.json();

        // Organize data by rooms
        const roomsMap = {};
        rawData.forEach((item) => {
          if (!roomsMap[item.room]) {
            roomsMap[item.room] = {};
          }

          // Store only the first occurrence of an appliance
          if (!roomsMap[item.room][item.appliance]) {
            roomsMap[item.room][item.appliance] = {
              name: item.appliance,
              status: item.status === "On",
              energyConsumption: item.energy_consumption,
              date: item.date,
              time: item.time,
            };
          }
        });

        // Convert room data to required structure
        const formattedRooms = Object.entries(roomsMap).reduce((acc, [room, appliances]) => {
          acc[room] = Object.values(appliances); // Get only unique appliances
          return acc;
        }, {});

        setRooms(Object.keys(formattedRooms));
        setAppliances(formattedRooms);
        setSelectedRoom(Object.keys(formattedRooms)[0]); // Default room
      } catch (error) {
        console.error("Error fetching data:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const handleRoomChange = (key) => {
    setSelectedRoom(key);
  };

  return (
    <>
      <Typography.Title style={{ fontSize: "18px", fontFamily: "Rubik" }}>Rooms</Typography.Title>
      {loading ? (
        <Spin size="large" style={{ display: "block", margin: "auto" }} />
      ) : (
        <Row gutter={[16, 16]}>
          {rooms.map((room) => (
            <Col xs={24} sm={12} md={8} lg={6} key={room}>
              <Card title={room} bordered>
                <p>{appliances[room]?.length} Appliances</p>
              </Card>
            </Col>
          ))}
        </Row>
      )}

      <div style={{ marginTop: 24 }}>
        <Typography.Title style={{ fontSize: "18px", marginBottom: 8, fontFamily: "Rubik" }}>
          Appliances
        </Typography.Title>
        {loading ? (
          <Spin size="large" style={{ display: "block", margin: "auto" }} />
        ) : (
          <Tabs
            activeKey={selectedRoom}
            onChange={handleRoomChange}
            type="line"
            tabBarGutter={16}
          >
            {rooms.map((room) => (
              <TabPane tab={room} key={room}>
                <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
                  {appliances[room]?.map((appliance, index) => (
                    <Col xs={24} sm={12} md={8} lg={6} key={index}>
                      <Card>
                        <p>
                          <strong>{appliance.name}</strong>
                        </p>
                        <p>Energy: {appliance.energyConsumption} kWh</p>
                        
                        <Switch defaultChecked={appliance.status} />
                      </Card>
                    </Col>
                  ))}
                </Row>
              </TabPane>
            ))}
          </Tabs>
        )}
      </div>
    </>
  );
}

export default Rooms;
