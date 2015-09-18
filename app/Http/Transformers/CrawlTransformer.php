<?php

namespace App\Http\Transformers;

use App\Article;
use App\Author;
use App\Crawl;
use League\Fractal\TransformerAbstract;

class CrawlTransformer extends TransformerAbstract
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

        if ($crawlable instanceof Author) {
            return $this->item($crawlable, new AuthorTransformer());
        }

        if ($crawlable instanceof Article) {
            return $this->item($crawlable, new ArticleTransformer());
        }
    }
}
