import React, {useEffect, useState} from "react";
import {PageHeader, Table, Input, Button, notification, Search} from "antd";
import {useDispatch, useSelector} from "react-redux";
import Column from "antd/lib/table/Column";
import DeployedBranchColumn from "./components/DeployedBranchColumn";
import NameColumn from "./components/NameColumn";
import LastDeploymentTimeColumn from "./components/LastDeploymentTimeColumn";
import DeploymentsTable from "./components/DeploymentsTable";
import {fetchDeployments} from '../../reducers/deploymentsReducer';
import {fetchJobs} from "../../reducers/jobsReducer";
import {LOAD_ERROR, LOAD_INIT, LOAD_LOADING} from "../../reducers/states";
import gitlabApi from "../../api/gitlabApi";
import {SyncOutlined} from "@ant-design/icons";

import './EvironmentPage.css';

export default (props) => {
    const envName = props.match.params.name;
    const env = useSelector(state => state.environments.items[envName]);
    const gitLabBaseURL = useSelector(state => state.config.data.gitLabBaseURL);
    const [branchPrefix, changeBranchPrefix] = useState('');
    const [isRequestLoading, changeRequestLoading] = useState(false);
    const {
        items: deployment,
        state: deploymentState,
        error: deploymentError
    } = useSelector(state => state.deployments);
    const userLinkTemplate = useSelector(state => state.config.data.userLinkTemplate)
    const projects = env ? env.projects : [];
    const dispatch = useDispatch();

    const renderDeployHistory = record => {
        if (!deployment || deploymentState === LOAD_INIT || deploymentState === LOAD_LOADING) {
            return (
                <div style={{display: 'flex', justifyContent: 'center'}}>
                    <SyncOutlined spin />
                </div>
            );
        }
        const deployHistory = deployment[record.id] && deployment[record.id].deployments || [];

        return (
            <DeploymentsTable
                deployments={deployHistory}
                userLinkTemplate={userLinkTemplate}
                gitLabBaseUrl={gitLabBaseURL}
                projectUrl={record.webURL}
            />
        )
    }

    useEffect(() => {
        if (projects.length > 0) {
            dispatch(fetchDeployments({envName, id: projects.map(({id}) => id)}));
        }
    }, [projects, envName]);

    const updateDeploymentsHistory = () => {
        dispatch(fetchDeployments({envName, id: projects.map(({id}) => id)}));
    }

    const onDeployByPrefix = (evt) => {
        evt.preventDefault();
        if (isRequestLoading) {
            return false;
        }

        changeRequestLoading(true);
        gitlabApi.deployByPrefix(envName, branchPrefix)
            .then(res => {
                projects.forEach(({id}) => {
                    dispatch(fetchJobs({envName, id}));
                })
                notification['success']({
                    message: 'Multi-deploy',
                    description: `Deploy for branches with '${branchPrefix}' prefix run`,
                });
            })
            .catch(err => {
                notification['error']({
                    message: 'Error occurred',
                    description: err.response.data.error,
                });
            })
            .finally(() => {
                changeRequestLoading(false)
            })
    }

    const aa = (asd) => {
        debugger
    }

    return <>
        <div className="env-page__header">
            <PageHeader
                className="site-page-header"
                onBack={() => window.location.pathname = "/"}
                title={envName}
            />
            <form onSubmit={onDeployByPrefix} className="env-page__deploy-all">
                <Input
                    disabled={isRequestLoading}
                    placeholder="Input prefix *"
                    value={branchPrefix}
                    onChange={({target: {value}}) => {
                        changeBranchPrefix(value);
                    }}
                />
                <Button
                    type="primary"
                    onClick={onDeployByPrefix}
                    loading={isRequestLoading}
                >
                    Deploy by prefix
                </Button>
            </form>
        </div>
        <Table dataSource={projects}
               pagination={false}
               loading={env === undefined}
               // expandRowByClick={true}
               expandable={{
                   expandedRowRender: renderDeployHistory,
               }}
               rowKey="id">
            <Column
                title="Name"
                render={project => <NameColumn project={project}/>}
            />
            <Column
                title="Deployed branch"
                render={project => (
                    <DeployedBranchColumn
                        id={project.id}
                        lastDeployment={project.lastDeployment}
                        envName={env && env.name}
                        updateDeploymentsHistory={updateDeploymentsHistory}
                    />
                )}
            />
            <Column
                title="Last deployment"
                render={project => (
                    <LastDeploymentTimeColumn
                        id={project.id}
                        lastDeployment={project.lastDeployment}
                        projectUrl={project.webURL}
                    />
                )}
                responsive={["md"]}
            />
            {/*<Column*/}
            {/*    title="User"*/}
            {/*    dataKey="User"*/}
            {/*    render={project => (*/}
            {/*        <UserColumn*/}
            {/*            user={project.lastDeployment ? project.lastDeployment.user : undefined}*/}
            {/*            gitLabBaseUrl={gitLabBaseURL}*/}
            {/*            linkTemplate={userLinkTemplate}*/}
            {/*        />*/}
            {/*    )}*/}
            {/*/>*/}
        </Table>
    </>
}
