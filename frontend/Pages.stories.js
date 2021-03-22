import React from "react";
import "antd/dist/antd.less";

import MainPage from "./Main.page";
import SignInPage from "./Signin.page";

export default {
    title: "Pages",
    component: MainPage,
};

const MainTmpl = (args) => <MainPage />;
const MainStory = MainTmpl.bind({});
MainStory.args = {};
MainStory.story = {
    name: "Main",
};

const SignInTmpl = (args) => <SignInPage />;
const SignInStory = SignInTmpl.bind({});
SignInStory.args = {};
SignInStory.story = {
    name: "SignIn",
};

export { MainStory, SignInStory };
