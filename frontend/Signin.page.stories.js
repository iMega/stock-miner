import React from "react";
import "antd/dist/antd.less";

import SignInPage from "./Signin.page";

export default {
    title: "Pages",
    component: SignInPage,
};

const SignInTmpl = (args) => <SignInPage />;
const SignInStory = SignInTmpl.bind({});
SignInStory.args = {};
SignInStory.story = {
    name: "SignIn",
};

export { SignInStory };
