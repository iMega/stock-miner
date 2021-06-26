import React from "react";
import { Menu, Layout } from "antd";
import {
    Switch,
    Route,
    Link,
    useLocation,
    BrowserRouter,
    StaticRouter,
} from "react-router-dom";

import { CAN_USE_DOM } from "./CanUseDom";
import Connector from "./Connector";
import PageStat from "./Stat";
import PageDealings from "./Dealings";
import { Add, List } from "./StockItem";
import Profile from "./Profile";

const { Sider, Content } = Layout;

const Main = () => (
    <Connector>
        <Layout>
            {CAN_USE_DOM ? (
                <BrowserRouter>
                    <Routing />
                </BrowserRouter>
            ) : (
                <StaticRouter>
                    <Routing />
                </StaticRouter>
            )}
        </Layout>
    </Connector>
);

const Routing = () => (
    <React.Fragment>
        <Sider breakpoint="lg" collapsedWidth="0">
            <MainMenu />
        </Sider>
        <Layout>
            <Content
                style={{
                    minHeight: "100vh",
                    background: "#fff",
                }}
            >
                <Switch>
                    {builderRoutes(mainMenuItems)}
                    {/* корневой должен быть последним, иначе роут не пашет */}
                    {buildRootRoute(mainMenuItems)}
                </Switch>
            </Content>
        </Layout>
    </React.Fragment>
);

const MainMenu = () => {
    let location = useLocation();
    return (
        <Menu
            defaultSelectedKeys={[location.pathname]}
            mode="inline"
            theme="dark"
        >
            {builderMenu(mainMenuItems)}
        </Menu>
    );
};

const mainMenuItems = [
    { path: "/profile", title: "Profile", page: <Profile /> },
    { path: "/", title: "Statistic", page: <PageStat /> },
    { path: "/dealings", title: "Dealings", page: <PageDealings /> },
    {
        path: "stock-item",
        title: "Stock item",
        sub: [
            {
                path: "/stock-item/list",
                title: "List",
                page: <List />,
            },
            {
                path: "/stock-item/add",
                title: "Add",
                page: <Add />,
            },
        ],
    },
];

const builderMenu = (menu) =>
    menu.map((item) =>
        hasProperty(item, "sub") ? (
            <Menu.SubMenu key={item.path} title={item.title}>
                {builderMenu(item.sub)}
            </Menu.SubMenu>
        ) : (
            builderItemMenu(item)
        )
    );

const builderItemMenu = (item) => (
    <Menu.Item key={item.path}>
        <Link to={item.path}>{item.title}</Link>
    </Menu.Item>
);

const builderRoutes = (menu) => {
    let r = [];
    menu.forEach(
        (i) =>
            i.path !== "/" &&
            (hasProperty(i, "sub")
                ? r.push(...builderRoutes(i.sub))
                : r.push(buildRoute(i)))
    );

    return r;
};

const buildRootRoute = (menu) =>
    menu.map((item) => item.path === "/" && buildRoute(item));

const buildRoute = (item) => (
    <Route key={item.path} path={item.path}>
        {item.page}
    </Route>
);

const hasProperty = (o, p) => Object.prototype.hasOwnProperty.call(o, p);

export default Main;
