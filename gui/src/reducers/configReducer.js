import {createSlice, createAsyncThunk} from "@reduxjs/toolkit";
import gitlabApi from "../api/gitlabApi";
import {LOAD_INIT, LOAD_SUCCESS, LOAD_LOADING, LOAD_ERROR} from "./states";


export const fetchConfig = createAsyncThunk(
    'config/fetch',
    async () => await gitlabApi.fetchConfig(),
)

const initialState = {
    data: {
        gitLabBaseURL: '',
        userLinkTemplate: '',
        gitLabAppId: '',
        oAuthEnabled: true,
        user: null,
    },
    state: LOAD_INIT,
}

export default createSlice({
        name: 'config',
        initialState,
        extraReducers: {
            [fetchConfig.pending]: (state) => ({
                data: state.data,
                state: LOAD_LOADING,
            }),
            [fetchConfig.rejected]: (state) => ({
                data: state.data,
                state: LOAD_ERROR
            }),
            [fetchConfig.fulfilled]: (state, action) => ({
                data: action.payload,
                state: LOAD_SUCCESS,
            }),

        }
    }
)
