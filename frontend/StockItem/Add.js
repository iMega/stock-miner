import React from "react";
import { PageHeader, Row, Col, Table, Input, Space, Button } from "antd";
import { useQuery, useMutation, gql } from "@apollo/client";

import Message from "../Message";
import Link from "../StockItemLink";

const MarketStockItemsND = gql`
    query MarketStockItems {
        marketStockItems {
            ticker
            figi
            isin
            minPriceIncrement
            lot
            currency
            name
        }
    }
`;

const Add = () => {
    const { loading, data } = useQuery(MarketStockItemsND);
    const [ds, setDs] = React.useState([]);
    const [selectedRowKeys, setSelectedRowKeys] = React.useState([]);
    const [filteredValue, setFilteredValue] = React.useState([]);

    React.useEffect(() => {
        if (loading === false && data) {
            setDs(data.marketStockItems);
        }
    }, [loading, data]);

    const [AddStockItem, { loading: muL }] = useMutation(AddStockItemAprovedND);

    const rowSelection = {
        selectedRowKeys,
        hideSelectAll: true,
        onChange: (keys) => setSelectedRowKeys(keys),
    };

    const onSearch = (value) => {
        const val = value.split(",");
        setFilteredValue(val.map((i) => i.trim().toLowerCase()));
    };

    const AddStockItemHandler = async (e) => {
        e.preventDefault();

        const items = ds.filter((i) => selectedRowKeys.includes(i.figi));
        const stockItems = items.map((i) => ({
            ticker: i.ticker,
            figi: i.figi,
            amountLimit: 0,
            transactionLimit: 0,
            currency: i.currency,
            startTime: 11,
            endTime: 20,
            maxPrice: 0,
            active: false,
        }));

        try {
            const { data } = await AddStockItem({
                variables: { in: stockItems },
            });
            if (data?.addStockItemApproved === true) {
                Message.Success();
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
            title="Avaliable stock items"
            extra={[
                <Button
                    key="1"
                    type={"primary"}
                    loading={muL}
                    disabled={selectedRowKeys.length === 0}
                    onClick={AddStockItemHandler}
                >
                    Add items
                </Button>,
            ]}
        >
            <Space direction="vertical" size="large">
                <Row>
                    <Col xs={24} lg={12}>
                        <Input.Search
                            addonBefore="Ticker:"
                            placeholder="AA, CAT, pdco"
                            allowClear
                            size="large"
                            onSearch={onSearch}
                        />
                    </Col>
                </Row>
                <Row>
                    <Col xs={24} lg={24}>
                        <Table
                            rowKey="figi"
                            rowSelection={rowSelection}
                            columns={columns(filteredValue)}
                            dataSource={ds}
                        />
                    </Col>
                </Row>
            </Space>
        </PageHeader>
    );
};

const columns = (filteredValue) => [
    {
        title: "Ticker",
        dataIndex: "ticker",
        key: "ticker",
        render: Link,
        filterMultiple: false,
        filteredValue: filteredValue,
        onFilter: (value, record) =>
            record.ticker.toLowerCase().includes(value),
    },
    {
        title: "Currency",
        dataIndex: "currency",
        key: "currency",
        responsive: ["xxl", "xl", "lg", "md"],
    },
    {
        title: "Name",
        dataIndex: "name",
        key: "name",
        responsive: ["xxl", "xl", "lg"],
    },
    {
        title: "FIGI",
        dataIndex: "figi",
        key: "figi",
    },
    {
        title: "ISIN",
        dataIndex: "isin",
        key: "isin",
        responsive: ["xxl", "xl", "lg"],
    },
    {
        title: "Lot",
        dataIndex: "lot",
        key: "lot",
        responsive: ["xxl", "xl", "lg", "md"],
    },
    {
        title: "Min price incr.",
        dataIndex: "minPriceIncrement",
        key: "minPriceIncrement",
        responsive: ["xxl", "xl", "lg", "md"],
    },
];

const AddStockItemAprovedND = gql`
    mutation AddStockItemAproved($in: [StockItemInput!]!) {
        addStockItemApproved(items: $in)
    }
`;

export default Add;
