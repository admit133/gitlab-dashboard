import React, {useState, useEffect} from "react";
import {PageHeader, Table, Input} from "antd";
import {useSelector} from "react-redux";
import Column from "antd/lib/table/Column";
import {LOAD_LOADING, LOAD_SUCCESS} from "../../reducers/states";
import {Link} from "react-router-dom";

export default () => {
    const envs = useSelector(state => state.environments);
    const items = envs.state === LOAD_SUCCESS ? Object.values(envs.items) : [];
    const [searchValue, changeSearchValue] = useState('');
    const [filteredValues, changeFilteredValues] = useState(undefined);

    const onChangeSearchValue = ({target: {value}}) => {
        if (value.trim() === '') {
            changeFilteredValues([])
        }
        changeSearchValue(value);
        const filtred = items.reduce((acc, curValue) => {
            if (curValue.name.toLowerCase().indexOf(value.toLowerCase()) !== -1) {
                acc.push(curValue);
            }
            return acc;
        }, []);
        changeFilteredValues(filtred);
    }

    const getItems = () => {
        if (filteredValues === undefined) {
            return items;
        } else if (Array.isArray(filteredValues)) {
            return filteredValues;
        }
    }

    return (
        <div>
            <div style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center'}}>
                <PageHeader
                    className="site-page-header"
                    title="Environments"
                />
                <Input
                    style={{width: '300px'}}
                    placeholder="Search"
                    allowClear
                    value={searchValue}
                    onChange={onChangeSearchValue}
                    autoFocus={true}
                />
            </div>
            <Table
                dataSource={getItems()}
               pagination={false}
               loading={envs.state === LOAD_LOADING}
               rowKey="Name"
            >
                <Column
                    title="Name"
                    render={(env) => {
                        return (
                            <Link
                                style={{display: 'inline-block', width: '100%', height: '100%'}}
                                to={`/environments/${env.name}`}
                            >
                                {env.name}
                            </Link>
                        )
                    }}
                />
            </Table>
        </div>
    )
}
