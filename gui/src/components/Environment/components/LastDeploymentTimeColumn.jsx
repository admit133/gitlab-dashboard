import React, {useEffect, useState} from "react";

import {Space} from "antd";
import moment from "moment";

import {LOAD_INIT, LOAD_LOADING, LOAD_SUCCESS} from "../../../reducers/states";
import {useSelector} from "react-redux";
import {SyncOutlined} from "@ant-design/icons";

export default props => {
    if (!props.lastDeployment) {
        return <span>no deployment in past</span>;
    }
    const {items: jobs, state: jobsState, error: jobsError} = useSelector(state => {
        return {
            ...state.jobs,
            items: state.jobs.items[props.id] ? state.jobs.items[props.id] : []
        };
    });

    const pipelineID = props.lastDeployment.deployable.pipeline.id;
    const projectUrl = props.projectUrl;
    const updatedAt = moment(props.lastDeployment.updatedAt)

    if (jobsState === LOAD_INIT) {
        return <SyncOutlined spin />
    }

    if ((jobs === null ||  jobs && jobs.length === 0) && jobsState === LOAD_SUCCESS) {
        return (
            <Space size="middle">
                <a href={`${projectUrl}/pipelines/${pipelineID}`}>{updatedAt.calendar()}</a>
            </Space>
        )
    }

    if (jobs && (jobs.status === 'created' || jobs.status === 'running' || jobs.status === 'pending')) {
        return (
            <Space size="middle">
                <a target="_brank" href={jobs.web_url}>Running...</a>
            </Space>
        )
    }

    if (!jobs || (Array.isArray(jobs) && jobs.length === 0)) {
        return null;
    }

    return (
        <Space size="middle">
            <a target="_brank" href={`${projectUrl}/pipelines/${jobs && jobs.pipeline.id}`}>{moment(jobs.finished_at).calendar()}</a>
        </Space>
    );
}
