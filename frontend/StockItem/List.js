import React from "react";
import { useQuery, useMutation, gql } from "@apollo/client";
import { PageHeader, Row, Col, Table, InputNumber, Space, Switch } from "antd";

import Link from "../StockItemLink";
import Message from "../Message";

const List = () => {
    let ds = [];

    const { loading, data } = useQuery(StockItemApprovedND, {
        fetchPolicy: "network-only",
    });

    if (loading === false && data) {
        ds = data.stockItemApproved;
    }

    const [updateStockItem] = useMutation(UpdateStockItemAprovedND);

    const settingsStockItemHandler = async (records) => {
        console.log(records);
        const stockItems = records.map((i) => ({
            ticker: i.ticker,
            figi: i.figi,
            amountLimit: 0,
            transactionLimit: 0,
            currency: i.currency,
            startTime: i.startTime,
            endTime: i.endTime,
            active: i.active,
            maxPrice: i.maxPrice,
        }));

        try {
            const { data } = await updateStockItem({
                variables: { in: stockItems },
            });
            if (data?.updateStockItemApproved === true) {
                return;
            }
            Message.Failure();
        } catch (e) {
            Message.Failure();
        }
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
                        columns={columns(settingsStockItemHandler, loading)}
                        dataSource={ds}
                        rowKey="figi"
                    />
                </Col>
            </Row>
        </PageHeader>
    );
};

const columns = (settingsStockItemHandler, loading) => [
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
                        settingsStockItemHandler([r]);
                    }}
                />
                -
                <HourInput
                    val={r.endTime}
                    onChange={(v) => {
                        r.endTime = v;
                        settingsStockItemHandler([r]);
                    }}
                />
            </Space>
        ),
    },
    {
        title: "Max purchase price",
        dataIndex: "maxPrice",
        key: "maxPrice",
        render: (_, r) => (
            <PriceInput
                val={r.maxPrice}
                onChange={(v) => {
                    r.maxPrice = v;
                    settingsStockItemHandler([r]);
                }}
            />
        ),
    },
    {
        title: "Active",
        dataIndex: "active",
        key: "active",
        render: (_, r) => (
            <Switch
                loading={loading}
                onChange={(v) => {
                    r.active = v;
                    settingsStockItemHandler([r]);
                }}
                checked={r.active}
            />
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

const PriceInput = ({ val, onChange }) => (
    <InputNumber
        min={0}
        defaultValue={val}
        style={{ width: "8em" }}
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
            currency
            startTime
            endTime
            active
            maxPrice
        }
    }
`;

const UpdateStockItemAprovedND = gql`
    mutation UpdateStockItemAproved($in: [StockItemInput!]!) {
        updateStockItemApproved(items: $in)
    }
`;

export default List;
