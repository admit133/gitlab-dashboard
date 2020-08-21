import {createSlice, createAsyncThunk} from "@reduxjs/toolkit";
import {LOAD_INIT, LOAD_SUCCESS, LOAD_LOADING, LOAD_ERROR} from "./states";
import gitlabApi from "../api/gitlabApi";

const initialState = {
    items: {},
    state: LOAD_INIT
}

export const fetchBranches = createAsyncThunk(
    'branches/fetch',
    async ({envName, id}) => await gitlabApi.fetchBranches(envName, id),
)

export default createSlice({
    name: 'branches',
    initialState,
    extraReducers: {
        [fetchBranches.fulfilled]: (state, action) => {
            const st = {
                ...state.items,
                [action.payload.projectId]: action.payload.branches
            }
            return {items: st, state: LOAD_SUCCESS}
        },
        [fetchBranches.pending]: () => ({items: {}, state: LOAD_LOADING}),
        [fetchBranches.rejected]: (state, action) => ({items: {}, state: LOAD_ERROR, error: action.error}),
    }
});
