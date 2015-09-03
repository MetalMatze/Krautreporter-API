<?php

namespace App\Http\Transformers;

use App\Crawl;

class CrawlTransformer
{
    protected $availableIncludes = [
        'crawlable'
    ];

    protected $defaultIncludes = [
        'crawlable'
    ];

    public function transform(Crawl $crawl)
    {
        return [
            'id' => (int)$crawl->id,
            'next_crawl' => $crawl->next_crawl->format(\DateTime::ISO8601),
        ];
    }

    public function includeCrawlable(Crawl $crawl)
    {
        $crawlable = $crawl->crawlable;

        return $this->item($crawlable, new AuthorTransformer());
    }
}
