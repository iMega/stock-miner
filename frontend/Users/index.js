import React from "react";
import { PageHeader, Row, Col, Table, Form, Input, Button } from "antd";
import { useQuery, useMutation, gql } from "@apollo/client";

import Message from "../Message";

const UsersND = gql`
    query Users {
        users {
            email
            name
            avatar
        }
    }
`;

const Page = () => {
    let ds = [];
    const { loading, data } = useQuery(UsersND);
    if (loading === false && data) {
        ds = data.users;
    }

    const [addUser] = useMutation(createUserND);
    const [removeUser] = useMutation(removeUserND);

    const onAdd = async (all) => {
        try {
            const { data } = await addUser({
                variables: { in: { email: all.email } },
            });
            if (data.createUser === true) {
                Message.Success();
                return;
            }
            Message.Failure();
        } catch (e) {
            Message.Failure();
        }
    };

    const onRemove = async (email) => {
        try {
            const { data } = await removeUser({
                variables: { in: { email: email } },
            });
            if (data.removeUser === true) {
                Message.Success();
                return;
            }
            Message.Failure();
        } catch (e) {
            Message.Failure();
        }
    };

    const columns = [
        {
            title: "Email",
            dataIndex: "email",
            key: "email",
        },
        {
            title: "Name",
            dataIndex: "name",
            key: "name",
        },
        {
            title: "Actions",
            render: (_, r) => (
                <Button type="danger" onClick={() => onRemove(r.email)}>
                    Remove
                </Button>
            ),
        },
    ];

    const [form] = Form.useForm();

    return (
        <React.Fragment>
            <PageHeader
                className="site-page-header"
                ghost={false}
                title="Users"
            >
                <Row>
                    <Col xs={24} lg={12}>
                        <Form
                            {...layout}
                            form={form}
                            onFinish={onAdd}
                            layout="inline"
                        >
                            <Form.Item
                                label="Email"
                                name="email"
                                htmlFor="email"
                                rules={[
                                    {
                                        required: true,
                                        message: "Please input email new user!",
                                    },
                                ]}
                            >
                                <Input
                                    type="email"
                                    id="email"
                                    style={{ width: "120%" }}
                                />
                            </Form.Item>
                            <Button
                                type="primary"
                                htmlType="submit"
                                style={{ width: "20%" }}
                            >
                                Add
                            </Button>
                        </Form>
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

const layout = {
    labelCol: { span: 6 },
    wrapperCol: { span: 16 },
};

const createUserND = gql`
    mutation CreateUser($in: UserInput!) {
        createUser(user: $in)
    }
`;

const removeUserND = gql`
    mutation RemoveUser($in: UserInput!) {
        removeUser(user: $in)
    }
`;

export default Page;
