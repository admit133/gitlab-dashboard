import React from "react";

import {Avatar, Space} from "antd";
import {useSelector} from "react-redux";
import {LOAD_INIT} from "../../../reducers/states";
import {SyncOutlined} from "@ant-design/icons";


const makeLink = (linkTemplate, username, gitLabBaseUrl) => {
    if (linkTemplate === "") {
        return `${gitLabBaseUrl}/${username}`
    }

    return linkTemplate.replace('{username}', username)
}

export default (props) => {
    if (props.user === undefined) {
        return <span/>;
    }
    const {items: jobs, state: jobsState, error: jobsError} = useSelector(state => state.jobs);

    const user = props.user;
    const linkTemplate = props.linkTemplate;
    const gitLabBaseUrl = props.gitLabBaseUrl;

    if (jobsState == LOAD_INIT) {
        return <SyncOutlined spin />
    }

    if (!jobs) {
        return (
            <Space size="middle">
                <Avatar icon={<img src={user.avatarURL} alt={`${user.username} Avatar`}/>} size={50}/>
                <a href={makeLink(linkTemplate, user.username, gitLabBaseUrl)}
                   target="_blank"
                   rel="noopener noreferrer"
                >@{user.username}
                </a>
            </Space>
        )
    }

    return (
        <Space size="middle">
            <Avatar icon={<img src={jobs.user.avatar_url} alt={`${jobs.user.username} Avatar`}/>} size={50}/>
            <a href={makeLink(linkTemplate, jobs.user.username, gitLabBaseUrl)}
               target="_blank"
               rel="noopener noreferrer"
            >@{jobs.user.username}
            </a>
        </Space>
    )
}
