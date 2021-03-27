import React from "react";
import { PageHeader, Tag, Row, Col, Statistic, Button } from "antd";
import { ArrowUpOutlined } from "@ant-design/icons";

const PageStat = () => (
    <PageHeader
        className="site-page-header"
        ghost={false}
        title="Статистика работы"
        subTitle="Статус бота:"
        tags={<Tag color="green">Работает</Tag>}
        extra={[
            <Button key="1" type="danger">
                Выключить
            </Button>,
        ]}
    >
        <Row>
            <Col xs={12} lg={3}>
                <Statistic title="Status" value="Pending" />
            </Col>
            <Col xs={12} lg={3}>
                <Statistic title="Price" prefix="$" value={568.08} />
            </Col>
            <Col xs={12} lg={3}>
                <Statistic title="Balance" prefix="$" value={3345.08} />
            </Col>
            <Col xs={12} lg={3}>
                <Statistic
                    title="Active"
                    value={11.28}
                    precision={2}
                    valueStyle={{ color: "#3f8600" }}
                    prefix={<ArrowUpOutlined />}
                    suffix="%"
                />
            </Col>
        </Row>
    </PageHeader>
);

export default PageStat;
