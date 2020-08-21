import environmentSlice from "./environmentsReducer"
import configSlice from "./configReducer"
import branchesReducer from "./branchesReducer"
import jobsReducer from "./jobsReducer"
import deploymentsReducer from "./deploymentsReducer"

export default {
    environments: environmentSlice.reducer,
    config: configSlice.reducer,
    branches: branchesReducer.reducer,
    jobs: jobsReducer.reducer,
    deployments: deploymentsReducer.reducer,
};
