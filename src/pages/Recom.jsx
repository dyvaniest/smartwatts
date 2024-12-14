import React, { useState, useEffect } from "react";
import {
  Layout,
  Row,
  Col,
  Card,
  Typography,
  Button,
  Input,
  Spin,
  Pagination,
  FloatButton,
  Modal,
} from "antd";
import { SearchOutlined } from "@ant-design/icons";

const { Title, Text } = Typography;
const { Search } = Input;

const Recommendations = () => {
  const [articles, setArticles] = useState([]);
  const [loading, setLoading] = useState(false);
  const [keyword, setKeyword] = useState("save energy");
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(6); // Number of articles per page
  const [isChatVisible, setIsChatVisible] = useState(false); // Kontrol popup

  const fetchArticles = async (query) => {
    setLoading(true);
    try {
      const response = await fetch(
        `https://newsapi.org/v2/everything?q=${query}&searchIn=title,description&language=en&apiKey=2889e3a78ba74d5e91948b45682dc6b8`
      );
      const data = await response.json();
      setArticles(data.articles || []);
    } catch (error) {
      console.error("Error fetching articles:", error);
    }
    setLoading(false);
  };

  useEffect(() => {
    fetchArticles(keyword); // Fetch articles on initial load
  }, []);

  const handleSearch = (value) => {
    if (value.trim() !== "") {
      setKeyword(value);
      setCurrentPage(1); // Reset to first page after new search
      fetchArticles(value);
    }
  };

  const handlePageChange = (page) => {
    setCurrentPage(page);
  };

  // Paginated data
  const startIndex = (currentPage - 1) * pageSize;
  const paginatedArticles = articles.slice(startIndex, startIndex + pageSize);

  return (
    <Layout.Content style={{ padding: "10px", margin: "0 16px" }}>

      {/* Konten Artikel */}
      <Row gutter={[16, 16]} style={{ marginBottom: 20 }}>
        <Col span={24}>
          <Title level={3} style={{ fontFamily: "Rubik", fontSize: "18px" }}>
            Energy Saving Recommendations
          </Title>
        </Col>
        <Col span={24}>
          <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
            <Col>
              <Search
                placeholder="Search for another articles..."
                enterButton={<SearchOutlined />}
                size="large"
                onSearch={handleSearch}
                style={{ width: 300 }}
              />
            </Col>
            <Col>
              <Pagination
                current={currentPage}
                pageSize={pageSize}
                total={articles.length}
                onChange={handlePageChange}
                showSizeChanger={false}
                simple
              />
            </Col>
          </Row>
        </Col>
        {loading ? (
          <Col span={24} style={{ textAlign: "center", marginTop: 20 }}>
            <Spin size="large" />
          </Col>
        ) : paginatedArticles.length > 0 ? (
          <>
            {paginatedArticles.map((article, index) => (
              <Col span={8} key={index}>
                <Card
                  hoverable
                  cover={
                    <img
                      alt={article.title}
                      src={article.urlToImage || "https://via.placeholder.com/300"}
                      style={{ height: 180, objectFit: "cover" }}
                    />
                  }
                >
                  <Title level={5} style={{ fontFamily: "Rubik" }}>
                    {article.title}
                  </Title>
                  <Text style={{ fontSize: "12px", color: "gray" }}>
                    {new Date(article.publishedAt).toLocaleDateString()}
                  </Text>
                  <Text style={{ display: "block", marginTop: 8 }}>
                    {article.description
                      ? `${article.description.substring(0, 100)}...`
                      : "No description available."}
                  </Text>
                  <Button
                    type="link"
                    href={article.url}
                    target="_blank"
                    style={{ marginTop: 8, padding: 0 }}
                  >
                    Read More
                  </Button>
                </Card>
              </Col>
            ))}
          </>
        ) : (
          <Col span={24}>
            <Text>No articles found for "{keyword}".</Text>
          </Col>
        )}
      </Row>
    </Layout.Content>
  );
};

export default Recommendations;
