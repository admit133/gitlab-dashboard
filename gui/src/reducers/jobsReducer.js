import {createSlice, createAsyncThunk} from "@reduxjs/toolkit";
import {LOAD_INIT, LOAD_SUCCESS, LOAD_LOADING, LOAD_ERROR} from "./states";
import gitlabApi from "../api/gitlabApi";

const initialState = {
    items: {},
    state: LOAD_INIT
}

export const fetchJobs = createAsyncThunk(
    'jobs/fetch',
    async ({envName, id}) => await gitlabApi.fetchJobs(envName, id),
)

export default createSlice({
    name: 'jobs',
    initialState,
    extraReducers: {
        [fetchJobs.fulfilled]: (state, action) => {
            const st = {
                ...state.items,
                [action.payload.projectId]: action.payload.job
            }
            return {items: st, state: LOAD_SUCCESS}
        },
        [fetchJobs.pending]: (state) => ({items: state.items, state: LOAD_LOADING}),
        [fetchJobs.rejected]: (state, action) => ({items: {}, state: LOAD_ERROR, error: action.error}),
    }
});
