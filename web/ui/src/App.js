import React, {Component} from 'react';
import 'react-bulma-components/dist/react-bulma-components.min.css';

import TenantList from './components/TenantList'
import Home from './components/Home'
import TopNavigation from './components/TopNavigation'
import {Route} from 'react-router-dom'

class App extends Component {


    render() {
        return (
            <div className="App">
                <TopNavigation/>

                <Route exact path="/" component={Home}/>
                <Route path="/tenants" component={TenantList}/>
            </div>
        );
    }
}

export default App;
