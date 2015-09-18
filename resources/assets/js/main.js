import React from 'react'
import request from 'superagent'
import {Table} from 'react-bootstrap'

class Status extends React.Component {

    constructor() {
        super();

        this.state = {
            crawls: []
        }
    }

    componentDidMount() {
        request
            .get('/crawls')
            .end((error, response) => {
                if (error == null) {
                    this.setState({
                        crawls: JSON.parse(response.text).data
                    })
                }
            })
    }

    render() {
        var crawlsRows = this.state.crawls.map((crawl) => {
            return (
                <tr>
                    <td>{crawl.id}</td>
                    <td>{crawl.next_crawl}</td>
                    <td>asdf</td>
                    <td>
                        <button className="btn btn-xs btn-info">more</button>
                    </td>
                </tr>
            )
        });

        return (
            <Table hover>
                <thead>
                <tr>
                    <th>Type</th>
                    <th>URL</th>
                    <th>Next Crawl</th>
                    <th></th>
                </tr>
                </thead>
                <tbody>
                {crawlsRows}
                </tbody>
            </Table>
        )
    }
}

React.render(<Status/>, document.querySelector('#status'));
