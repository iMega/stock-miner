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
import Message from "../Message";
import useToggle from "../UseToggle";

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

const MiningStopND = gql`
    mutation MiningStop {
        disableStockItemsApproved
    }
`;

const MiningStartND = gql`
    mutation MiningStart {
        enableStockItemsApproved
    }
`;

const StateND = gql`
    query Slots {
        user {
            role
        }
        settings {
            miningStatus
        }
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
    const [stopMining] = useMutation(MiningStopND);
    const [startMining] = useMutation(MiningStartND);

    const [isGlobalMining, toggleGlobalMining] = useToggle();
    const [stopGlobalMining] = useMutation(GlobalMiningStopND);
    const [startGlobalMining] = useMutation(GlobalMiningStartND);

    let ds = {
        user: {
            role: "",
        },
        settings: {
            miningStatus: false,
        },
        slots: [],
    };
    const { loading, data } = useQuery(StateND);
    if (loading === false) {
        ds = data;
        setStatusMining(ds.settings.miningStatus);
    }

    const switchMining = async () => {
        try {
            const { data } = statusMining
                ? await stopMining()
                : await startMining();
            if (
                data.enableStockItemsApproved === true ||
                data.disableStockItemsApproved === true
            ) {
                Message.Success();
                setStatusMining((v) => !v);
                return;
            }
            Message.Failure();
        } catch (e) {
            Message.Failure();
        }
    };

    const switchGlobalMining = async () => {
        try {
            const { data } = isGlobalMining
                ? await stopGlobalMining()
                : await startGlobalMining();
            if (
                data.globalMiningStop === true ||
                data.globalMiningStart === true
            ) {
                Message.Success();
                toggleGlobalMining();
                return;
            }
            Message.Failure();
        } catch (e) {
            Message.Failure();
        }
    };

    const buttonsBar = [
        <Button
            key="1"
            type={statusMining ? "danger" : "primary"}
            onClick={switchMining}
        >
            Switch {statusMining ? "OFF" : "ON"}
        </Button>,
    ];

    if (ds?.user.role === "root") {
        buttonsBar.push(
            <Button key="2" danger onClick={switchGlobalMining}>
                Switch Global {isGlobalMining ? "OFF" : "ON"}
            </Button>
        );
    }

    const [tagColor, tagTitle] = status(ds?.settings.miningStatus);

    return (
        <React.Fragment>
            <PageHeader
                className="site-page-header"
                ghost={false}
                title="Statistic"
                subTitle="Status:"
                tags={<Tag color={tagColor}>{tagTitle}</Tag>}
                extra={buttonsBar}
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
                    <Table columns={columns} dataSource={ds?.slots} />
                </Col>
            </Row>
        </React.Fragment>
    );
};

const status = (v) => (v === true ? ["green", "Runnung"] : ["red", "Stopped"]);

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
