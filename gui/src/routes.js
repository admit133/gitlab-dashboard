import React from 'react';
import App from "./App";
import {Route, Switch} from "react-router";
import EnvironmentPage from './components/Environment/EnvironmentPage';
import EnvironmentsPage from './components/Environments/EnvironmentsPage';
import NotFound from './components/NotFound/NotFound';

export default () => (
    <App>
        <Switch>
            <Route exact path="/" component={EnvironmentsPage}/>
            <Route path="/environments/:name" component={EnvironmentPage}/>
            <Route path="*" component={NotFound}/>
        </Switch>
    </App>
)
