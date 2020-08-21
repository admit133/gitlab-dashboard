import React from "react";

import {Avatar, Space} from "antd";
import GitlabOutlined from "@ant-design/icons/lib/icons/GitlabOutlined";

export default (props) => {
    const project = props.project;
    return <Space size="middle">
        <Avatar
            icon={<GitlabOutlined />}
            src={project.avatarURL} alt={`${project.nameWithNamespace} Avatar`} size={"large"}/>
        <a target="_brank"  href={project.webURL}>{project.nameWithNamespace}</a>
    </Space>
}
