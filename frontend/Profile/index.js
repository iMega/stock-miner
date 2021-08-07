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
import RulePrice from "./RulePrice";

const MarketCredentialsND = gql`
    mutation MarketCredentials($in: MarketCredentialsInput!) {
        marketCredentials(creds: $in)
    }
`;

const slotND = gql`
    mutation Slot($in: SlotSettingsInput!) {
        slot(global: $in)
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
            marketCommission
            grossMargin
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
    const [creds, setCreds] = React.useState();
    const [slot, setSlot] = React.useState();
    const { loading, data } = useQuery(ProfileAndSettingsND);
    const [saveCreds] = useMutation(MarketCredentialsND);
    const [saveSlot] = useMutation(slotND);
    React.useEffect(() => {
        if (loading === false && data) {
            if (
                data.settings.marketCredentials !== null &&
                data.settings.marketCredentials.length > 0
            ) {
                console.log(
                    "+++++++++++",
                    data,
                    data.settings.marketCredentials.length > 0
                );
                setCreds({
                    provider: data.settings.marketCredentials[0].name,
                    url: data.settings.marketCredentials[0].apiUrl,
                    token: data.settings.marketCredentials[0].token,
                });
            }
            setSlot({
                volume: data.settings.slot.volume,
            });
        }
    }, [loading, data]);

    const onSaveCreds = async (all) => {
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

    const onSaveSlot = async (all) => {
        try {
            const { data } = await saveSlot({
                variables: {
                    in: {
                        volume: all.volume,
                    },
                },
            });
            if (data.slot === true) {
                Message.Success();
                return;
            }
            Message.Failure();
        } catch (e) {
            Message.Failure();
        }
    };

    const [credForm] = Form.useForm();
    const selectProvider = (e) => {
        const creds = data.settings.marketCredentials;
        if (creds === null) {
            return;
        }

        const idx = creds.findIndex((i) => i.name === e.target.value);
        credForm.setFieldsValue({
            provider: creds[idx].name,
            url: creds[idx].apiUrl,
            token: creds[idx].token,
        });
    };
    const [slotForm] = Form.useForm();
    const Settings = () => (
        <React.Fragment>
            <Form
                {...layout}
                form={credForm}
                onFinish={onSaveCreds}
                initialValues={creds}
            >
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
                    <Radio.Group onChange={selectProvider}>
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
            <Form
                {...layout}
                form={slotForm}
                onFinish={onSaveSlot}
                initialValues={slot}
            >
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
                    <InputNumber min={0} max={10} />
                </Form.Item>
                <Form.Item {...tailLayout}>
                    <Button type="primary" htmlType="submit">
                        Save
                    </Button>
                </Form.Item>
            </Form>
            <Divider />
            <RulePrice {...data} />
        </React.Fragment>
    );

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
                <Button key="1" type={"danger"} href="/logout">
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
