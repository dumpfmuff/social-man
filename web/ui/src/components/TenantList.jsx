import React, {Component} from 'react';

import axios from "axios";

import Columns from "react-bulma-components/lib/components/columns";
import Card from "react-bulma-components/lib/components/card";
import Media from "react-bulma-components/lib/components/media";
import Image from "react-bulma-components/lib/components/image";
import Heading from "react-bulma-components/lib/components/heading";
import Content from "react-bulma-components/lib/components/content";

export default class TenantList extends Component {

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
            <Columns>
                {this.renderTenants()}
            </Columns>
        );
    }

    renderTenants() {
        let key = 0
        return this.state.tenants.map((tenant) =>
            <Columns.Column size={2} key={key++}>{TenantList.renderTenantCard(tenant)}</Columns.Column>
        );
    }

    static renderTenantCard(tenant) {
        console.log(tenant)
        return (<Card>
            <Card.Image size="4by3" src="https://bulma.io/images/placeholders/640x480.png"/>
            <Card.Content>
                <Media>
                    <Media.Item renderas="figure" position="left">
                        <Image renderas="p" size={64} alt="64x64"
                               src="http://bulma.io/images/placeholders/128x128.png"/>
                    </Media.Item>
                    <Media.Item>
                        <Heading size={4}>{tenant.lastname}, {tenant.surname}</Heading>
                        <Heading subtitle size={6}>
                            @johnsmith
                        </Heading>
                    </Media.Item>
                </Media>
                <Content>
                    joined:
                </Content>
            </Card.Content>
        </Card>);
    }
}