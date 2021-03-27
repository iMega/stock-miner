import React from "react";
import { PageHeader, Row, Col, Table } from "antd";

const List = () => {
    const columns = [
        {
            title: "Ticker",
            dataIndex: "ticker",
            key: "ticker",
            render: (text) => <a>{text}</a>,
        },
        {
            title: "FIGI",
            dataIndex: "figi",
            key: "figi",
        },
        {
            title: "Limit",
            dataIndex: "limit",
            key: "limit",
        },
    ];

    const data = [
        {
            key: "1",
            ticker: "AAPL",
            figi: "32",
        },
    ];

    return (
        <PageHeader
            className="site-page-header"
            ghost={false}
            title="List of stock items"
        >
            <Row>
                <Col xs={24} lg={12}>
                    <Table columns={columns} dataSource={data} />
                </Col>
            </Row>
        </PageHeader>
    );
};

export default List;
