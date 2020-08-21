import React from 'react';
import {Button, Layout} from "antd";import {
    GitlabOutlined,
} from '@ant-design/icons';

import './Login.css';

const {Content} = Layout;

export default ({gitLabBaseURL, gitLabBaseAppId}) => (
    <Layout className="layout login-page">
        <Content className="login-page__content">
            <div className="login-page__logo" />
            <Button
                danger
                className="login-page__login-btn"
                size="large"
                type="primary"
                icon={<GitlabOutlined />}
                href={`${gitLabBaseURL}/oauth/authorize?client_id=${gitLabBaseAppId}&redirect_uri=${window.location.origin}/oauth/code&response_type=code&scope=read_user`}
            >
                Sign in with GitLab
            </Button>
        </Content>
    </Layout>
);
