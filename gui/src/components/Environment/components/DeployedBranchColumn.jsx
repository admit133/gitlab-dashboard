import React, {useState, useEffect}  from "react";
import {useDispatch, useSelector} from "react-redux";
import gitlabApi from '../../../api/gitlabApi';

import {Space, Select, Spin, Alert, notification} from "antd";
import {
    SyncOutlined,
    CheckCircleTwoTone,
    CloseCircleTwoTone,
    ReloadOutlined
} from '@ant-design/icons';

import {LOAD_ERROR, LOAD_INIT, LOAD_LOADING, LOAD_SUCCESS} from "../../../reducers/states";
import {fetchBranches} from "../../../reducers/branchesReducer";
import {fetchJobs} from "../../../reducers/jobsReducer";

const { Option } = Select;

export default ({lastDeployment, id, envName, updateDeploymentsHistory}) => {
    const {items: branches, state: branchesState, error: branchesError} = useSelector(state => {
        return {
            ...state.branches,
            items: state.branches.items[id] ? state.branches.items[id] : []
        };
    });
    const {items: jobs, state: jobsState, error: jobsError} = useSelector(state => {
        return {
            ...state.jobs,
            items: state.jobs.items[id] ? state.jobs.items[id] : []
        };
    });

    const [ref, changeRef] = useState('');
    const [selectedBranch, changeSelectedBranch] = useState('');
    const [timeoutId, changeTimeoutId] = useState(null);
    const [isLoading, changeIsLoading] = useState(false);
    const [alert, setAlert] = useState(null);
    const dispatch = useDispatch();

    useEffect(() => {
        dispatch(fetchJobs({envName, id}));
        dispatch(fetchBranches({envName, id}));
    }, [envName, id]);

    useEffect(() => {
        dispatch(fetchJobs({envName, id}));
        dispatch(fetchBranches({envName, id}));
    }, []);

    const requestJob = async (envName, id) => {
        const {job} = await gitlabApi.fetchJobs(envName, id);
        if (job && (job.status === 'pending' || job.status === 'running' || job.status === 'created')) {
            setTimeout(() => {
                requestJob(envName, id)
            }, 6000)
        } else {
            updateDeploymentsHistory();
            dispatch(fetchJobs({envName, id}));
        }
    }

    useEffect(() => {
        const init = () => {
            changeRef(jobs.ref)
            changeSelectedBranch(jobs.ref)
            if (jobs.status === 'pending' || jobs.status === 'running' || jobs.status === 'created') {
                //clearInterval(timeoutId);
                changeIsLoading(true);
                /*const intervalId = setInterval(() => {
                    dispatch(fetchJobs({envName, id}));
                }, 5000);
                changeTimeoutId(intervalId);*/
                requestJob(envName, id);
            } else {
                changeIsLoading(false);
            }
        }

        if (Array.isArray(jobs) && jobs.length > 0) {
            init();
        } else if (jobs && !Array.isArray(jobs)) {
            init();
        } else {
            changeRef( lastDeployment && lastDeployment.ref || '')
            changeSelectedBranch(lastDeployment && lastDeployment.ref || '')
            changeIsLoading(false);
        }

        return () => {
            if (timeoutId) {
                clearInterval(timeoutId)
                changeTimeoutId(null);
            }
        }

    }, [jobs, envName, id]);

    useEffect(() => {
        if (timeoutId !== null) {
            if (jobs !== null) {
                if (jobs.status !== 'pending' && jobs.status !== 'running' && jobs.status !== 'created') {
                    clearInterval(timeoutId);
                    changeIsLoading(false);
                }
            }
        }
    }, [timeoutId, jobs]);

    const onChange = branchName => {
        const branchInfo = branches.find(branch => branch.name === branchName);
        changeSelectedBranch(branchName);
        changeIsLoading(true);

        gitlabApi.deployBranch(envName, id, branchName, branchInfo.commit.id)
            .then(() => {
                dispatch(fetchJobs({envName, id}));
                setTimeout(() => {
                    clearInterval(timeoutId);
                    const intervalId = setInterval(() => {
                        dispatch(fetchJobs({envName, id}));
                    }, 5000);
                    changeTimeoutId(intervalId);
                }, 1500);
            })
            .catch(err => {
                changeIsLoading(false);
                changeSelectedBranch(ref);
                notification['error']({
                    message: 'Error occurred',
                    description: err.response.data.error,
                });
            })
    }

    const onRebuildBranch = () => {
        onChange(selectedBranch);
    }

    const shouldDisableSelect = () => {
        return jobs && (jobs.status === 'pending' || jobs.status === 'running' || jobs.status === 'created');
    }

    const renderIcon = () => {
        if (isLoading) {
            return (
                <span title={jobs && jobs.status || ''}><SyncOutlined spin /></span>
            )
        }

        if (jobs && jobs.length === 0) {
            if (!lastDeployment) {
                return null;
            }
            const lastDeployStatus = lastDeployment.deployable.pipeline.status
            switch (lastDeployStatus) {
                case 'success':
                    return (
                        <a href="">
                            <CheckCircleTwoTone twoToneColor="#52c41a"/>
                        </a>
                    )
                case 'canceled':
                case 'skipped':
                case 'failed':
                    return (
                        <CloseCircleTwoTone twoToneColor="#eb2f96" />
                    )
            }
        }

        if (jobs && jobs.status === 'success') {
            return (
                <a target="_brank" href={jobs.web_url} title="success">
                    <CheckCircleTwoTone twoToneColor="#52c41a"/>
                </a>
            )
        }

        if (jobs && (jobs.status === 'canceled' || jobs.status === 'skipped' || jobs.status === 'failed')) {
            return (
                <a target="_brank" href={jobs.web_url} title={jobs && jobs.status || ''}>
                    <CloseCircleTwoTone twoToneColor="#eb2f96" />
                </a>
            )
        }

        return null
    }

    if (branchesState === LOAD_ERROR && branchesError && branchesError.message) {
        setAlert(<Alert message={branchesError.message} banner closable/>)
    }

    const renderSelect = () => {
        return (
            <Select
                showSearch
                style={{ width: 200 }}
                placeholder="Select a person"
                optionFilterProp="children"
                notFoundContent={branchesState === LOAD_LOADING ? (
                    <Spin size="small" />
                ) : null}
                disabled={isLoading || shouldDisableSelect()}
                onChange={onChange}
                defaultValue={ref}
                value={selectedBranch}
                filterOption={(input, option) =>
                    option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
                }
            >
                {
                    branches.map(branch => (
                        <Option key={branch.name} value={branch.name}>{branch.name}</Option>
                    ))
                }
            </Select>
        )
    }

    const renderRedeployButton = () => {
        if (selectedBranch.trim() === '') {
            return null;
        }
        const style = isLoading ? {
            opacity: 0.5,
            cursor: 'not-allowed',
        } : {
            cursor: 'pointer',
        }

        if (branchesState === LOAD_SUCCESS) {
            return isLoading ? (
                <span title="Redeploy current branch" style={style}>
                    <ReloadOutlined />
                </span>
            ) : (
                <span onClick={onRebuildBranch}  title="Redeploy current branch" style={style}>
                    <ReloadOutlined />
                </span>
            )
        }

        return null;
    }

    return (
        <>
            <Space size="middle">
                { renderSelect() }
                { renderRedeployButton() }
                { renderIcon() }
            </Space>
            { alert }
        </>
    );
}
