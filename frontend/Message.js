import { message } from "antd";

const Success = () => {
    message.success("The operation was completed successfully.");
};

const Failure = () => {
    message.error("Please try again later.");
};

const Message = {
    Success,
    Failure,
};

export default Message;
