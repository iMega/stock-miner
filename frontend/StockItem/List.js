import React from "react";
import { useQuery, gql } from "@apollo/client";
import { PageHeader, Row, Col, Table, InputNumber, Space } from "antd";

import Link from "../StockItemLink";

const List = () => {
    let ds = [
        {
            ticker: "AAPL",
            figi: "BBG000B9XRY4",
            amountLimit: 0,
            transactionLimit: 0,
            startTime: 10,
            endTime: 22,
        },
        {
            ticker: "AMZN",
            figi: "BBG000BVPV84",
            amountLimit: 0,
            transactionLimit: 0,
            startTime: 12,
            endTime: 20,
        },
    ];

    const { loading, data } = useQuery(StockItemApprovedND, {
        fetchPolicy: "network-only",
    });

    if (loading === false && data) {
        ds = data.stockItemApproved;
    }

    const settingsStockItemHandler = (rec) => {
        console.log(rec);
    };

    return (
        <PageHeader
            className="site-page-header"
            ghost={false}
            title="List of stock items"
        >
            <Row>
                <Col xs={24} lg={12}>
                    <Table
                        loading={loading}
                        columns={columns(settingsStockItemHandler)}
                        dataSource={ds}
                        rowKey="figi"
                    />
                </Col>
            </Row>
        </PageHeader>
    );
};

const columns = (settingsStockItemHandler) => [
    {
        title: "Ticker",
        dataIndex: "ticker",
        key: "ticker",
        render: Link,
    },
    {
        title: "FIGI",
        dataIndex: "figi",
        key: "figi",
    },
    {
        title: "Work hours",
        key: "work-hours",
        render: (_, r) => (
            <Space>
                <HourInput
                    val={r.startTime}
                    onChange={(v) => {
                        r.startTime = v;
                        settingsStockItemHandler(r);
                    }}
                />
                <HourInput
                    val={r.endTime}
                    onChange={(v) => {
                        r.endTime = v;
                        settingsStockItemHandler(r);
                    }}
                />
            </Space>
        ),
    },
];

const HourInput = ({ val, onChange }) => (
    <InputNumber
        min={0}
        max={23}
        defaultValue={val}
        style={{ width: "4em" }}
        onChange={onChange}
    />
);

const StockItemApprovedND = gql`
    query StockItemApproved {
        stockItemApproved {
            ticker
            figi
            amountLimit
            transactionLimit
        }
    }
`;

export default List;
