import React, {Component} from 'react';
import axios from 'axios';

class App extends Component {

    constructor(props) {
        super(props)
        this.state = {tenants: []};
        this.renderTenants = this.renderTenants.bind(this)
    }


    componentWillMount() {
        let url = '/api/tenants'
        let that = this;
        axios.get(url).then(response => {
            console.log(response);
            that.setState({
                tenants: response.data
            })
        });
    }

    render() {
        return (
            <div className="App">
                <header className="App-header">
                    <h1 className="App-title">Welcome to Social Man</h1>
                </header>
                <ul>{this.renderTenants()}</ul>
            </div>
        );
    }

    renderTenants() {
        let key = 0
        return this.state.tenants.map((tenant) =>
            <li key={key++}>{tenant.lastname}, {tenant.surname}</li>)
    }

}

export default App;
