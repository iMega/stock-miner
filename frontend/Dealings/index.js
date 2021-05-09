import React from "react";
import { PageHeader, Row, Col, Space, Statistic, Table } from "antd";
import { useQuery, gql } from "@apollo/client";

import StockItemLink from "../StockItemLink";

const DealingsND = gql`
    query Dealings {
        dealings {
            id
            ticker
            figi
            startPrice
            changePrice
            buyingPrice
            targetPrice

            qty
            profit
            amountSpent
            amountIncome
            totalProfit
            salePrice

            buyAt
            duration
            sellAt

            currency
        }
    }
`;

const Page = () => {
    let ds = [];
    const { loading, data } = useQuery(DealingsND);
    if (loading === false && data) {
        ds = data?.dealings;
    }

    return (
        <React.Fragment>
            <PageHeader
                className="site-page-header"
                ghost={false}
                title="Dealings"
            >
                <Row>
                    <Col xs={24} sm={12} lg={12}>
                        <Space direction="horizontal" size="large" wrap={true}>
                            <Statistic title="Profit" prefix="$" value={0} />
                            <Statistic title="Portfolio" prefix="$" value={0} />
                            <Statistic title="Expected" prefix="$" value={0} />
                        </Space>
                    </Col>
                </Row>
            </PageHeader>
            <Row>
                <Col span={22} offset={1} md={16} lg={22}>
                    <Table columns={columns} dataSource={ds} />
                </Col>
            </Row>
        </React.Fragment>
    );
};

const Locale = "ru-RU";

const Currency = (locale, currency, amount) =>
    new Intl.NumberFormat(locale, {
        style: "currency",
        currency: currency,
    }).format(amount);

const columns = [
    {
        title: "Sell time",
        dataIndex: "DateTime",
        key: "sellAt",
        render: (v, r) => new Date(r.sellAt).toLocaleString(Locale),
    },
    {
        title: "Duration",
        dataIndex: "duration",
        key: "duration",
        align: "right",
        render: (v) => new Date(75600000 + v * 1000).toLocaleTimeString(Locale),
    },
    {
        title: "Ticker",
        dataIndex: "ticker",
        key: "ticker",
        render: StockItemLink,
    },
    {
        title: "Buy",
        dataIndex: "buyingPrice",
        key: "buyingPrice",
        align: "right",
        render: (v, r) => Currency(Locale, r.currency, v),
    },
    {
        title: "Sale",
        dataIndex: "salePrice",
        key: "salePrice",
        align: "right",
        render: (v, r) => Currency(Locale, r.currency, v),
    },
    {
        title: "Profit",
        dataIndex: "profit",
        key: "profit",
        responsive: ["xxl", "xl", "lg", "md"],
        align: "right",
        render: (v, r) => Currency(Locale, r.currency, v),
    },
    {
        title: "Qty",
        dataIndex: "qty",
        key: "qty",
        responsive: ["xxl", "xl", "lg", "md"],
        align: "right",
    },
    {
        title: "Spent Amount",
        dataIndex: "amountSpent",
        key: "amountSpent",
        responsive: ["xxl", "xl", "lg", "md"],
        align: "right",
        render: (v, r) => Currency(Locale, r.currency, v),
    },
    {
        title: "Income Amount",
        dataIndex: "amountIncome",
        key: "amountIncome",
        responsive: ["xxl", "xl", "lg", "md"],
        align: "right",
        render: (v, r) => Currency(Locale, r.currency, v),
    },
    {
        title: "Total Profit",
        dataIndex: "totalProfit",
        key: "totalProfit",
        responsive: ["xxl", "xl", "lg", "md"],
        align: "right",
        render: (v, r) => Currency(Locale, r.currency, v),
    },
];

export default Page;
