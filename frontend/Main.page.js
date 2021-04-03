import React from "react";
import { Menu, Layout } from "antd";
import { BrowserRouter as Router, Switch, Route, Link } from "react-router-dom";

import Connector from "./Connector";
import PageStat from "./Stat";
import { Add, List } from "./StockItem";
import Profile from "./Profile";

const { Sider, Content } = Layout;

const LinkStat = "/";
const LinkProfile = "/profile";
const LinkStockItemList = "/stock-item/list";
const LinkStockItemAdd = "/stock-item/add";

const Main = () => (
    <Connector>
        <Layout>
            <Router basename="/">
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
                            <Route path={LinkStat}>
                                <PageStat />
                            </Route>
                        </Switch>
                    </Content>
                </Layout>
            </Router>
        </Layout>
    </Connector>
);

const MainMenu = () => (
    <Menu defaultSelectedKeys={["2"]} mode="inline" theme="dark">
        <Menu.Item key="1">
            <Link to={LinkProfile}>Profile</Link>
        </Menu.Item>
        <Menu.Item key="2">
            <Link to={LinkStat}>Statistic</Link>
        </Menu.Item>
        <Menu.SubMenu key="sub1" title="Stock item">
            <Menu.Item key="sub1-1">
                <Link to={LinkStockItemList}>List</Link>
            </Menu.Item>
            <Menu.Item key="sub1-2">
                <Link to={LinkStockItemAdd}>Add</Link>
            </Menu.Item>
        </Menu.SubMenu>
    </Menu>
);

export default Main;
