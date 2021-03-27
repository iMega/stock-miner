import React from "react";
import { PageHeader, Row, Col, Table } from "antd";
import { useQuery, gql } from "@apollo/client";

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

    React.useEffect(() => {
        if (loading === false && data) {
            setDs(data.marketStockItems);
        }
    }, [loading, data]);

    return (
        <PageHeader
            className="site-page-header"
            ghost={false}
            title="Avaliable stock items"
        >
            <Row>
                <Col xs={24} lg={24}>
                    <Table rowKey="figi" columns={columns} dataSource={ds} />
                </Col>
            </Row>
        </PageHeader>
    );
};

const columns = [
    {
        title: "Ticker",
        dataIndex: "ticker",
        key: "ticker",
        render: (text) => <a>{text}</a>,
    },
    {
        title: "Currency",
        dataIndex: "currency",
        key: "currency",
    },
    {
        title: "Name",
        dataIndex: "name",
        key: "name",
        render: (text) => <a>{text}</a>,
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
    },
    {
        title: "Lot",
        dataIndex: "lot",
        key: "lot",
    },
    {
        title: "Min price incr.",
        dataIndex: "minPriceIncrement",
        key: "minPriceIncrement",
    },
];

export default Add;
