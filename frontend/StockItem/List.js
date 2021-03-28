import React from "react";
import { useQuery, gql } from "@apollo/client";
import { PageHeader, Row, Col, Table } from "antd";

import Link from "./Link";

const List = () => {
    const { loading, data } = useQuery(StockItemApprovedND);
    const [ds, setDs] = React.useState([]);

    React.useEffect(() => {
        if (loading === false && data) {
            setDs(data.stockItemApproved);
        }
    }, [loading, data]);

    const columns = [
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
            title: "Limit",
            dataIndex: "limit",
            key: "limit",
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
                    <Table columns={columns} dataSource={ds} />
                </Col>
            </Row>
        </PageHeader>
    );
};

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
