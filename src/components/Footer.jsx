import React from 'react';
// import footer from antd
import { Layout } from 'antd';
import "../index.css";

function Footer() {
    return (
        <Layout.Footer style={{ textAlign: 'left', background: "none", padding: "24px",}}>
            <p> <strong>SmartWatts: Smart Home Energy Management System</strong></p> 
            <p>©{new Date().getFullYear()} Created by Divany Pangestika</p>
        </Layout.Footer>
    );
    
}

export default Footer;