import React from "react";
import { Form, InputNumber, Button } from "antd";
import { useMutation, gql } from "@apollo/client";

import Message from "../Message";
import { Layout, TailLayout } from "./Style";

const Block = (props) => {
    const initialValues = {
        commission: props?.settings?.marketCommission || 0,
        margin: props?.settings?.grossMargin || 0,
    };
    const [save] = useMutation(rulePriceND);

    const [form] = Form.useForm();

    const onSave = async (all) => {
        try {
            const { data } = await save({
                variables: {
                    in: {
                        marketCommission: all.commission,
                        grossMargin: all.margin,
                    },
                },
            });
            if (data.rulePrice === true) {
                Message.Success();
                return;
            }
            Message.Failure();
        } catch (e) {
            Message.Failure();
        }
    };

    return (
        <Form
            {...Layout}
            form={form}
            onFinish={onSave}
            initialValues={initialValues}
        >
            <Form.Item
                label="Commission"
                name="commission"
                rules={[
                    {
                        required: true,
                        message: "Please input commission of market!",
                    },
                ]}
            >
                <InputNumber min={0} />
            </Form.Item>
            <Form.Item
                label="Gross margin"
                name="margin"
                rules={[
                    {
                        required: true,
                        message: "Please input gross margin!",
                    },
                ]}
            >
                <InputNumber min={0} />
            </Form.Item>
            <Form.Item {...TailLayout}>
                <Button type="primary" htmlType="submit">
                    Save
                </Button>
            </Form.Item>
        </Form>
    );
};

const rulePriceND = gql`
    mutation RulePrice($in: RulePriceInput!) {
        rulePrice(global: $in)
    }
`;

export default Block;
