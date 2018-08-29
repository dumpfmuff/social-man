import React, {Component} from 'react';
import {Link} from 'react-router-dom';

import Navbar from 'react-bulma-components/lib/components/navbar';

export default class TopNavigation extends Component {

    constructor(props) {
        super(props)
        this.state = {
            open: true
        }
    }

    render() {
        return (

            <Navbar
                // color={select('Color', colors)}
                fixed='top'
                active={true}
                transparent={false}
            >
                <Navbar.Brand>
                    <Navbar.Item renderas="a" href="#">
                        <img
                            src="https://bulma.io/images/bulma-logo.png"
                            alt="Bulma: a modern CSS framework based on Flexbox"
                            width="112"
                            height="28"
                        />
                    </Navbar.Item>
                    <Navbar.Burger
                        active={this.state.open}
                        onClick={() =>
                            this.setState(() => {
                                this.state.open = !this.state.open;
                            })
                        }
                    />
                </Navbar.Brand>
                <Navbar.Menu active={this.state.open}>
                    <Navbar.Container>
                        <Link to="/" className="navbar-item">Home</Link>
                        <Link to="/tenants" className="navbar-item">Tenants</Link>
                    </Navbar.Container>
                </Navbar.Menu>
            </Navbar>

        );
    }
}