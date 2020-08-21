import axios from 'axios';

const baseUrl = '';

const fetchRequest = async url => {
    const response = await axios.get(url);
    const resJson = response.data;
    if (resJson.error) {
        throw new Error(resJson.error);
    }
    return resJson;
}

export default class {
    static async fetchEnvironments() {
        const { environments } = await fetchRequest(`${baseUrl}/environments`);
        return {
            environments: environments.sort((a, b) => a.name > b.name),
        }
    }
    static async fetchConfig() {
       return await fetchRequest(`${baseUrl}/config`);
    }
    static async fetchBranches(envName, projectId) {
        const  { branches } = await fetchRequest(`${baseUrl}/environments/${envName}/projects/${projectId}/repository/branches`);
        return {
            projectId,
            branches,
        }
    }
    static async fetchJobs(envName, projectId) {
        const  { job } = await fetchRequest(`${baseUrl}/environments/${envName}/projects/${projectId}/jobs`);
        return {
            projectId,
            job,
        }
    }
    static async fetchDeployments(envName, projectIds) {
        const deployments = {};
        for (let i = 0; i < projectIds.length; i++) {
            const curDeployments = await fetchRequest(`${baseUrl}/environments/${envName}/projects/${projectIds[i]}/deployments`);
            deployments[projectIds[i]] = curDeployments;
        }

        return deployments;
    }
    static async deployBranch(envName, projectId, branchName, branchId) {
        return axios({
            method: 'post',
            url: `${baseUrl}/environments/${envName}/projects/${projectId}/jobs`,
            headers: {
                'Content-Type': 'application/json'
            },
            data: {
                ref: branchName,
                sha: branchId,
            }
        });
    }
    static deployByPrefix(envName, prefix) {
        return axios({
            method: 'post',
            url: `${baseUrl}/environments/${envName}/jobs`,
            headers: {
                'Content-Type': 'application/json'
            },
            data: {
                query: prefix,
            }
        });
    }
}
