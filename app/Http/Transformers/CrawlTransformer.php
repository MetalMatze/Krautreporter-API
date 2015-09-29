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
        if ($crawl->crawlable_type == Author::class) {
            $crawlable_type = 'author';
        } else {
            $crawlable_type = 'article';
        }

        return [
            'id' => (int)$crawl->id,
            'next_crawl' => $crawl->next_crawl->format(\DateTime::ISO8601),
            'crawlable_type' => $crawlable_type,
        ];
    }

    public function includeCrawlable(Crawl $crawl)
    {
        $crawlable = $crawl->crawlable;

        if ($crawlable instanceof Author) {
            return $this->item($crawlable, new AuthorTransformer());
        }

        if ($crawlable instanceof Article) {
            return $this->item($crawlable, new ArticleIncludeTransformer());
        }
    }
}
