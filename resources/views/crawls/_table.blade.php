<table class="table table-bordered">
    <thead>
    <tr>
        <th>URL</th>
        <th>Next Crawl</th>
    </tr>
    </thead>
    <tbody>
    @foreach($crawls as $crawl)
        <tr>
            <td>{{ $crawl->crawlable->url }}</td>
            <td>{{ $crawl->next_crawl }}</td>
        </tr>
    @endforeach
    </tbody>
</table>
