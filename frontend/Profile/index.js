import React from "react";
import {
    PageHeader,
    Row,
    Col,
    Form,
    Input,
    InputNumber,
    Button,
    Radio,
    Divider,
} from "antd";
import { useQuery, useMutation, gql } from "@apollo/client";

import Message from "../Message";

const MarketCredentialsND = gql`
    mutation MarketCredentials($in: MarketCredentialsInput!) {
        marketCredentials(creds: $in)
    }
`;

const ProfileAndSettingsND = gql`
    query Profile {
        user {
            email
            name
            avatar
        }
        settings {
            slot {
                volume
            }
            marketCredentials {
                name
                token
                apiUrl
            }
        }
    }
`;

const radioStyle = {
    display: "block",
    height: "30px",
    lineHeight: "30px",
};

const layout = {
    labelCol: { span: 6 },
    wrapperCol: { span: 16 },
};

const tailLayout = {
    wrapperCol: { offset: 6, span: 16 },
};

const Profile = () => {
    const [provider, setProvider] = React.useState(1);
    const { loading, data } = useQuery(ProfileAndSettingsND);
    const [saveCreds] = useMutation(MarketCredentialsND);
    React.useEffect(() => {
        if (loading === false && data) {
            console.log(data);
            // setDs(data.marketStockItems);
        }
    }, [loading, data]);

    const onFieldsChange2 = async (all) => {
        console.log(all);
        try {
            const { data } = await saveCreds({
                variables: {
                    in: {
                        name: all.provider,
                        apiUrl: all.url,
                        token: all.token,
                    },
                },
            });
            if (data.marketCredentials === true) {
                Message.Success();
                return;
            }
            Message.Failure();
        } catch (e) {
            Message.Failure();
        }
    };

    const [credForm] = Form.useForm();
    const [slotForm] = Form.useForm();
    const Settings = () => (
        <React.Fragment>
            <Form {...layout} form={credForm} onFinish={onFieldsChange2}>
                <Form.Item
                    label="Market provider"
                    name="provider"
                    rules={[
                        {
                            required: true,
                            message: "Please select url provider!",
                        },
                    ]}
                >
                    <Radio.Group
                        onChange={(e) => setProvider(e.target.value)}
                        value={provider}
                        defaultValue={provider}
                    >
                        <Radio value={"tinkoff"} style={radioStyle}>
                            Tinkoff production
                        </Radio>
                        <Radio value={"tinkoff-sandbox"} style={radioStyle}>
                            Tinkoff sandbox
                        </Radio>
                    </Radio.Group>
                </Form.Item>
                <Form.Item
                    label="API url"
                    name="url"
                    rules={[
                        {
                            required: true,
                            message: "Please input api url!",
                        },
                    ]}
                >
                    <Input type="url" placeholder="https://api.example.com" />
                </Form.Item>
                <Form.Item
                    label="Token"
                    name="token"
                    rules={[
                        {
                            required: true,
                            message: "Please input your token!",
                        },
                    ]}
                >
                    <Input.TextArea rows={4} />
                </Form.Item>
                <Form.Item {...tailLayout}>
                    <Button type="primary" htmlType="submit">
                        Save
                    </Button>
                </Form.Item>
            </Form>
            <Divider />
            <Form {...layout} form={slotForm} onValuesChange={onFieldsChange2}>
                <Form.Item
                    label="Volume of slot"
                    name="volume"
                    rules={[
                        {
                            required: true,
                            message: "Please input volume of slot!",
                        },
                    ]}
                >
                    <InputNumber min={1} max={10} />
                </Form.Item>
                <Form.Item {...tailLayout}>
                    <Button type="primary">Save</Button>
                </Form.Item>
            </Form>
        </React.Fragment>
    );

    console.log("RENDER");

    return (
        <PageHeader
            className="site-page-header"
            ghost={false}
            title={data?.user.name}
            avatar={{
                alt: "your avatar",
                size: "large",
                src: data?.user.avatar,
            }}
            extra={[
                <Button key="1" type={"danger"} href="/google/logout">
                    Logout
                </Button>,
            ]}
        >
            <Row>
                <Col xs={24} lg={12}>
                    <Settings />
                </Col>
            </Row>
        </PageHeader>
    );
};

export default Profile;
