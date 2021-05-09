import React from "react";
import { Menu, Layout } from "antd";
import { CAN_USE_DOM } from "./CanUseDom";
import { Switch, Route, Link } from "react-router-dom";
import { BrowserRouter, StaticRouter } from "react-router-dom";

import Connector from "./Connector";
import PageStat from "./Stat";
import PageDealings from "./Dealings";
import { Add, List } from "./StockItem";
import Profile from "./Profile";

const { Sider, Content } = Layout;

const LinkStat = "/";
const LinkProfile = "/profile";
const LinkDealings = "/dealings";
const LinkStockItemList = "/stock-item/list";
const LinkStockItemAdd = "/stock-item/add";

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
                    <Route path={LinkProfile}>
                        <Profile />
                    </Route>
                    <Route path={LinkStockItemList}>
                        <List />
                    </Route>
                    <Route path={LinkStockItemAdd}>
                        <Add />
                    </Route>
                    <Route path={LinkDealings}>
                        <PageDealings />
                    </Route>
                    {/* должен быть последним, иначе роут не пашет */}
                    <Route path={LinkStat}>
                        <PageStat />
                    </Route>
                </Switch>
            </Content>
        </Layout>
    </React.Fragment>
);

const MainMenu = () => (
    <Menu defaultSelectedKeys={["statistic"]} mode="inline" theme="dark">
        <Menu.Item key="profile">
            <Link to={LinkProfile}>Profile</Link>
        </Menu.Item>
        <Menu.Item key="statistic">
            <Link to={LinkStat}>Statistic</Link>
        </Menu.Item>
        <Menu.Item key="dealings">
            <Link to={LinkDealings}>Dealings</Link>
        </Menu.Item>
        <Menu.SubMenu key="stock-item" title="Stock item">
            <Menu.Item key="stock-item-list">
                <Link to={LinkStockItemList}>List</Link>
            </Menu.Item>
            <Menu.Item key="stock-item-add">
                <Link to={LinkStockItemAdd}>Add</Link>
            </Menu.Item>
        </Menu.SubMenu>
    </Menu>
);

export default Main;
