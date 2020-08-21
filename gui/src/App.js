import React from 'react';
import './App.css';
import {Layout, Button} from "antd";
import Header from "./components/Header/Header";
import Login from './components/Login/Login';
import {useSelector} from "react-redux";
import {LOAD_INIT, LOAD_SUCCESS, LOAD_LOADING, LOAD_ERROR} from "./reducers/states";
import {SyncOutlined} from "@ant-design/icons";

const {Content, Footer} = Layout;

function App(props) {
    const user = useSelector(state => state.config.data.user)
    const {
        data: {
            oAuthEnabled,
            gitLabBaseURL,
            gitLabAppId: gitLabBaseAppId
        },
        state: configState,
    } = useSelector(state => state.config)

    if (configState === LOAD_INIT || configState === LOAD_LOADING) {
        return  null;
    }

    if (oAuthEnabled && user === null) {
        return <Login gitLabBaseURL={gitLabBaseURL} gitLabBaseAppId={gitLabBaseAppId}/>
    }
    return (
        <Layout className="layout">
            <Header/>
            <Content style={{padding: '0 50px'}}>
                {props.children}
            </Content>
            <Footer style={{textAlign: 'center'}}>©{(new Date()).getFullYear()} GitLab Dashboard</Footer>
        </Layout>
    );
}

export default App;
