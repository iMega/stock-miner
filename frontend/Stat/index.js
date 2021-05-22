import React from "react";
import {
    PageHeader,
    Tag,
    Row,
    Col,
    Space,
    Statistic,
    Button,
    Table,
} from "antd";
import { useQuery, useMutation, gql } from "@apollo/client";

import StockItemLink from "../StockItemLink";

const MemStatsND = gql`
    subscription MemStats {
        memStats {
            alloc
            totalAlloc
            sys
        }
    }
`;

const GlobalMiningStopND = gql`
    mutation GlobalMiningStop {
        globalMiningStop
    }
`;

const GlobalMiningStartND = gql`
    mutation GlobalMiningStart {
        globalMiningStart
    }
`;

const SlotsND = gql`
    query Slots {
        slots {
            id
            ticker
            figi

            startPrice
            changePrice
            buyingPrice
            targetPrice
            profit

            qty
            amountSpent
            targetAmount
            totalProfit

            currency
            currentPrice
        }
    }
`;

const PageStat = () => {
    // const { loading, data } = useSubscription(MemStatsND, {
    //     shouldResubscribe: true,
    //     fetchPolicy: "network-only",
    //     onSubscriptionData: (data) => console.log("new data", data),
    // });

    const [statusMining, setStatusMining] = React.useState(false);
    const [stopMining, a] = useMutation(GlobalMiningStopND);
    const [startMining, b] = useMutation(GlobalMiningStartND);

    let ds = [];
    const { loading, data } = useQuery(SlotsND);
    if (loading === false) {
        ds = data?.slots;
    }

    const switchMining = () => {
        if (statusMining) {
            stopMining();
        } else {
            startMining();
        }
        setStatusMining(!statusMining);
    };

    return (
        <React.Fragment>
            <PageHeader
                className="site-page-header"
                ghost={false}
                title="Statistic"
                subTitle="Status:"
                tags={<Tag color="green">Running</Tag>}
                extra={[
                    <Button
                        key="1"
                        type={statusMining ? "danger" : "primary"}
                        onClick={switchMining}
                    >
                        Switch {statusMining ? "OFF" : "ON"}
                    </Button>,
                ]}
            >
                <Row>
                    <Col xs={24} sm={12} lg={12}>
                        <Space direction="horizontal" size="large" wrap={true}>
                            <Statistic title="Profit" prefix="$" value={0} />
                            <Statistic title="Portfolio" prefix="$" value={0} />
                            <Statistic title="Expected" prefix="$" value={0} />
                        </Space>
                    </Col>
                    <Col xs={0} sm={12} lg={12}>
                        <Space direction="horizontal" size="large" wrap={true}>
                            <Statistic title="Alloc" value={"1.102 MB"} />
                            <Statistic title="Total" value={"324.47 MB"} />
                            <Statistic title="Sys" value={"71.09 MB"} />
                        </Space>
                    </Col>
                </Row>
            </PageHeader>
            <Row>
                <Col span={22} offset={1} md={16} lg={20}>
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
        title: "Ticker",
        dataIndex: "ticker",
        key: "ticker",
        render: StockItemLink,
    },
    {
        title: "Start price",
        dataIndex: "startPrice",
        key: "startPrice",
        responsive: ["xxl", "xl", "lg", "md"],
        align: "right",
        render: (v, r) => Currency(Locale, r.currency, v),
    },
    {
        title: "Buy",
        dataIndex: "buyingPrice",
        key: "buyingPrice",
        align: "right",
        render: (v, r) => Currency(Locale, r.currency, v),
    },
    {
        title: "Current",
        dataIndex: "currentPrice",
        key: "currentPrice",
        align: "right",
        render: (v, r) => Currency(Locale, r.currency, v),
    },
    {
        title: "Target",
        dataIndex: "targetPrice",
        key: "targetPrice",
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
        title: "Target Amount",
        dataIndex: "targetAmount",
        key: "targetAmount",
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

export default PageStat;
