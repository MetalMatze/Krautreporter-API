import React from 'react'
import request from 'superagent'
import {Table} from 'react-bootstrap'
import moment from 'moment'

class Status extends React.Component {

    constructor() {
        super();

        this.state = {
            crawls: []
        };

        this.update = this.update.bind(this);
    }

    componentDidMount() {
        this.update();
        setInterval(this.update, 30 * 1000);
    }

    render() {
        var crawlsRows = this.state.crawls.map((crawl) => {
            return (
                <tr key={crawl.id}>
                    <td>{crawl.crawlable.data.id}</td>
                    <td>{crawl.crawlable_type == 'author' ? 'Author' : 'Article'}</td>
                    <td>{crawl.crawlable_type == 'author' ? crawl.crawlable.data.name : crawl.crawlable.data.title}</td>
                    <td>{moment(crawl.next_crawl, moment.ISO_8601).fromNow()}</td>
                </tr>
            )
        });

        return (
            <Table hover>
                <thead>
                <tr>
                    <th style={{width: 45}}>ID</th>
                    <th style={{width: 65}}>Type</th>
                    <th>URL</th>
                    <th style={{width: 150}}>Next Crawl</th>
                </tr>
                </thead>
                <tbody>
                {crawlsRows}
                </tbody>
            </Table>
        )
    }

    update() {
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
}

React.render(<Status/>, document.querySelector('#status'));
