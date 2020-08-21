import React from 'react';
import ReactDOM from 'react-dom';
import 'antd/dist/antd.css';
import './index.css';
import store from './app/store';
import {Provider} from 'react-redux';
import * as serviceWorker from './serviceWorker';
import {Router} from "react-router";
import {createBrowserHistory} from "history";
import Routes from "./routes"
import {fetchConfig} from "./reducers/configReducer";


store.dispatch(fetchConfig());

ReactDOM.render(
    <Provider store={store}>
        {/*https://github.com/ant-design/ant-design/issues/22493*/}
        {/*<React.StrictMode>*/}
        <Router history={createBrowserHistory()}>
            <Routes/>
        </Router>
        {/*</React.StrictMode>*/}
    </Provider>
    ,
    document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
