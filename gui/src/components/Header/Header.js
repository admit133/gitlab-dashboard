import {Link} from "react-router-dom";
import {Alert, Avatar, Layout, Menu, Space} from "antd";
import {SettingOutlined} from "@ant-design/icons";
import React from "react";
import {withRouter} from "react-router";
import SubMenu from "antd/es/menu/SubMenu";
import {useDispatch, useSelector} from "react-redux";

import {LOAD_ERROR, LOAD_INIT} from "../../reducers/states";
import {fetchEnvironments} from "../../reducers/environmentsReducer";

import './Header.css';

function Header(props) {
    const oAuthEnabled = useSelector(state => state.config.data.oAuthEnabled)
    const gitLabBaseURL = useSelector(state => state.config.data.gitLabBaseURL)
    const gitLabBaseAppId = useSelector(state => state.config.data.gitLabAppId)
    const user = useSelector(state => state.config.data.user)
    const environments = useSelector(state => state.environments);
    const dispatch = useDispatch();
    const {location} = props;
    let alert = null;

    if (environments.state === LOAD_INIT) {
        dispatch(fetchEnvironments())
    }

    if (environments.state === LOAD_ERROR) {
        alert = <Alert message={environments.error.message} banner closable/>
    }

    let loginMenu;

    if (oAuthEnabled) {
        if (user === null) {
            loginMenu =
                <div style={{float: 'right'}}>
                    <Menu.Item key={"login"}>
                        <a href={`${gitLabBaseURL}/oauth/authorize?client_id=${gitLabBaseAppId}&redirect_uri=${window.location.origin}/oauth/code&response_type=code&scope=read_user`}>Login</a>
                    </Menu.Item>
                </div>
        } else {
            const avatar = user ? (
                <div>
                    <Space size={"middle"}>
                        <Avatar icon={<img src={user.avatarURL} alt={`${user.username} Avatar`}/>} size={25}/>
                        {user.username}
                    </Space>
                </div>
            ) : null;
            loginMenu =
                <SubMenu style={{float: 'right'}} icon={avatar}>
                    <Menu.Item key="logout">
                        <a href={`/oauth/logout`}>Logout</a>
                    </Menu.Item>
                </SubMenu>
        }
    }

    return (
        <div>
            <Layout.Header style={{display: 'flex', alignItems: 'center'}}>
                <Link to="/" className="menu__logo">Gitlab Dashboard</Link>
                <Menu theme="dark" mode="horizontal" defaultSelectedKeys={[location.pathname]} style={{ width: '95%'}}>
                    <SubMenu icon={<SettingOutlined/>} title="Envs">
                        <Menu.Item key="All">
                            <Link to={`/`}>All</Link>
                        </Menu.Item>
                        {Object.keys(environments.items).map((envName) =>
                            <Menu.Item key={envName}>
                                <Link to={`/environments/${envName}`}>{envName}</Link>
                            </Menu.Item>
                        )}
                    </SubMenu>
                    {loginMenu}
                </Menu>
            </Layout.Header>
            {alert}
        </div>
    )
}


export default withRouter(Header)
