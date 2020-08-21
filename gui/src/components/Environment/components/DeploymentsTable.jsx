import React, {useEffect} from "react";
import moment from "moment";
import {Avatar, Space, Table} from "antd";
import Column from "antd/lib/table/Column";
import {CheckCircleTwoTone, CloseCircleTwoTone, SyncOutlined} from "@ant-design/icons";

const makeLink = (linkTemplate, username, gitLabBaseUrl) => {
    if (linkTemplate === "") {
        return `${gitLabBaseUrl}/${username}`
    }

    return linkTemplate.replace('{username}', username)
}

const renderStatus = status => {
    switch (status) {
        case 'success':
            return (
                <span title="success">
                    <CheckCircleTwoTone twoToneColor="#52c41a"/>
                </span>
            )
        case 'canceled':
        case 'skipped':
        case 'failed':
            return (
                <span title={status}>
                    <CloseCircleTwoTone twoToneColor="#eb2f96" />
                </span>
            )
        case 'pending':
        case 'runnning':
        case 'created':
            return (
                <span title={status}>
                    <SyncOutlined spin />
                </span>
            )
    }
};

const renderUser = (user, linkTemplate, gitLabBaseUrl) => (
    <Space size="middle">
        <Avatar icon={<img src={user.avatarURL} alt={`${user.username} Avatar`}/>} size={25}/>
        <a href={makeLink(linkTemplate, user.username, gitLabBaseUrl)}
           target="_blank"
           rel="noopener noreferrer"
        >
            @{user.username}
        </a>
    </Space>
);

export default ({deployments, userLinkTemplate, gitLabBaseUrl, projectUrl}) => (
    <Table dataSource={deployments}
           pagination={false}
           rowKey="id">
        <Column
            title="Branch"
            render={deploy => <a href={`${projectUrl}/-/tree/${deploy.ref}`} target="_blank"> {deploy.ref}</a>}
        />
        <Column
            title="Last update"
            render={deploy => <a href={`${projectUrl}/pipelines/${deploy.deployable.pipeline.id}`}
                                 target="_blank">{moment(deploy.updatedAt).calendar()}</a>}
        />
        <Column
            title="Status"
            render={deploy => renderStatus(
                deploy.deployable.pipeline.status,
            )}
        />
    </Table>
)
