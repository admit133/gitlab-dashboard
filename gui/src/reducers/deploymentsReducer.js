import {createSlice, createAsyncThunk} from "@reduxjs/toolkit";
import {LOAD_INIT, LOAD_SUCCESS, LOAD_LOADING, LOAD_ERROR} from "./states";
import gitlabApi from "../api/gitlabApi";

const initialState = {
    items: null,
    state: LOAD_INIT
}

export const fetchDeployments = createAsyncThunk(
    'deployments/fetch',
    async ({envName, id}) => await gitlabApi.fetchDeployments(envName, id),
)

export default createSlice({
    name: 'deployments',
    initialState,
    extraReducers: {
        [fetchDeployments.fulfilled]: (state, action) => ({items: action.payload, state: LOAD_SUCCESS}),
        [fetchDeployments.pending]: () => ({items: null, state: LOAD_LOADING}),
        [fetchDeployments.rejected]: (state, action) => ({items: null, state: LOAD_ERROR, error: action.error}),
    },
});
