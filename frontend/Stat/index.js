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
import { useMutation, gql } from "@apollo/client";

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

const PageStat = () => {
    // const { loading, data } = useSubscription(MemStatsND, {
    //     shouldResubscribe: true,
    //     fetchPolicy: "network-only",
    //     onSubscriptionData: (data) => console.log("new data", data),
    // });

    const [statusMining, setStatusMining] = React.useState(false);
    const [stopMining, a] = useMutation(GlobalMiningStopND);
    const [startMining, b] = useMutation(GlobalMiningStartND);

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
                title="Статистика работы"
                subTitle="Статус бота:"
                tags={<Tag color="green">Работает</Tag>}
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
                <Col span={22} offset={1} md={16} lg={12}>
                    <Table columns={columns} dataSource={ds} />
                </Col>
            </Row>
        </React.Fragment>
    );
};

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
    },
    {
        title: "Buy",
        dataIndex: "buy",
        key: "buy",
        align: "right",
    },
    {
        title: "Current",
        dataIndex: "current",
        key: "current",
        align: "right",
    },
    {
        title: "Target",
        dataIndex: "target",
        key: "target",
        align: "right",
    },
    {
        title: "Profit",
        dataIndex: "profit",
        key: "profit",
        responsive: ["xxl", "xl", "lg", "md"],
        align: "right",
    },
];

const ds = [
    {
        ticker: "AAPL",
        startPrice: 121.98,
        buy: 121.7,
        current: 120.25,
        target: 121.9,
        profit: 0.2,
    },
    {
        ticker: "AAPL",
        startPrice: 121.42,
        buy: 121.39,
        current: 121.39,
        target: 121.43,
        profit: 0.2,
    },
    {
        ticker: "PDCO",
        startPrice: 30.98,
        buy: 30.7,
        current: 30.75,
        target: 30.82,
        profit: 0.12,
    },
];

export default PageStat;
