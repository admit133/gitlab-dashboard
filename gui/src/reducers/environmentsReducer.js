import {createSlice, createAsyncThunk} from "@reduxjs/toolkit";
import {LOAD_INIT, LOAD_SUCCESS, LOAD_LOADING, LOAD_ERROR} from "./states";
import gitlabApi from "../api/gitlabApi";

const initialState = {
    items: {},
    state: LOAD_INIT
}

export const fetchEnvironments = createAsyncThunk(
    'environments/fetch',
    async () => await gitlabApi.fetchEnvironments(),
)

export default createSlice({
    name: 'environments',
    initialState,
    extraReducers: {
        [fetchEnvironments.fulfilled]: (state, action) => ({items: covertArrayToMapWithKey(action.payload.environments, 'name'), state: LOAD_SUCCESS}),
        [fetchEnvironments.pending]: () => ({items: {}, state: LOAD_LOADING}),
        [fetchEnvironments.rejected]: (state, action) => ({items: {}, state: LOAD_ERROR, error: action.error}),
    }
})

function covertArrayToMapWithKey(items, key) {
    let result = {}
    for (let i = 0; i < items.length; i++) {
        result[items[i][key]] = items[i]
    }

    return result;
}
