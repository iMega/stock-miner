import React from "react";
import "antd/dist/antd.less";

import MainPage from "./Main.page";

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

export { MainStory };
