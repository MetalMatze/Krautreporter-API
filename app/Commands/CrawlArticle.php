<?php namespace App\Commands;

use App\Article;
use App\Image;
use Carbon\Carbon;
use Goutte\Client;
use Illuminate\Contracts\Bus\SelfHandling;
use Illuminate\Contracts\Queue\ShouldBeQueued;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;
use Illuminate\Support\Facades\Log;
use Symfony\Component\DomCrawler\Crawler;

class CrawlArticle extends Command implements SelfHandling, ShouldBeQueued
{
    use InteractsWithQueue, SerializesModels;
    /**
     * @var Article
     */
    private $article;

    /**
     * Create a new command instance.
     *
     * @return void
     */
    public function __construct(Article $article)
    {
        $this->article = $article;
    }

    /**
     * Execute the command.
     *
     * @return void
     */
    public function handle()
    {
        $client = new Client();
        $client->followRedirects(false);

        $crawler = $client->request('GET', 'https://krautreporter.de' . $this->article->url);

        if ($client->getResponse()->getStatus() == 200) {
            if ($this->article->trashed()) {
                $this->article->restore();
            }

            $articleNode = $crawler->filter('main article.article.article--full')->eq(1);
            $articleHeaderNode = $articleNode->filter('header.article__header');
            $articleContentNode = $articleNode->filter('.article__content');

            $articleDateText = $articleHeaderNode->filter('h2.meta')->text();

            $this->article->headline = $articleHeaderNode->filter('h1.article__title')->text();
            $this->article->date = Carbon::createFromFormat('d.m.Y', $articleDateText);

            $articleImageNode = $articleHeaderNode->filter('.media__img img');

            if ($articleHeaderNode->filter('figure.media--fullwidth')->count() > 0) {
                $youtubeId = $articleHeaderNode->filter('figure.media--fullwidth div')->attr('data-video-id');

                $thumbnails = [
                    ['src' => "https://img.youtube.com/vi/$youtubeId/mqdefault.jpg", 'width' => 300],
                    ['src' => "https://img.youtube.com/vi/$youtubeId/sddefault.jpg", 'width' => 600],
                    ['src' => "https://img.youtube.com/vi/$youtubeId/maxresdefault.jpg", 'width' => 2000],
                ];

                foreach ($thumbnails as $thumbnail) {
                    $image = Image::where('imageable_type', '=', 'App\Article')
                        ->where('imageable_id', '=', $this->article->id)
                        ->where('src', '=', $thumbnail['src'])
                        ->first();

                    if ($image == null) {
                        $image = new Image();
                    }

                    $image->src = $thumbnail['src'];
                    $image->width = $thumbnail['width'];

                    $this->article->images()->save($image);
                }
            }

            if ($articleImageNode->count() > 0) {
                $widths = $articleImageNode->attr("srcset");
                preg_match("/(.*) 300w, (.*) 600w, (.*) 1000w, (.*) 2000w/", $widths, $matches);

                foreach ($matches as $index => $match) {
                    if ($index == 0) {
                        continue;
                    }

                    switch ($index) {
                        case 1:
                            $width = 300;
                            break;
                        case 2:
                            $width = 600;
                            break;
                        case 3:
                            $width = 1000;
                            break;
                        case 4:
                            $width = 2000;
                            break;
                    }

                    $image = Image::where('imageable_type', '=', 'App\Article')
                        ->where('imageable_id', '=', $this->article->id)
                        ->where('width', '=', $width)
                        ->first();

                    if ($image == null) {
                        $image = new Image();
                    }

                    $image->src = getenv("URL_KRAUTREPORTER") . $match;
                    $image->width = $width;

                    $this->article->images()->save($image);
                }
            }

            $this->article->excerpt = trim($articleContentNode->filter('h2.gamma')->text());

            $articleNode->filter('.article__content h2.gamma')->each(function (Crawler $crawler) {
                $node = $crawler->getNode(0);
                $node->parentNode->removeChild($node);
            });

            $this->article->content = trim($articleNode->filter('.article__content')->html());

            $this->article->save();

            $this->calculateNextCrawlDate();

        } else {
            $this->article->crawl->delete();
            $this->article->delete();
        }
    }

    private function calculateNextCrawlDate()
    {
        if (Carbon::now()->diffInMonths($this->article->date) > 0) {
            $nextCrawl = Carbon::now()->addDays(3);
        } elseif (Carbon::now()->diffInWeeks($this->article->date) > 0) {
            $nextCrawl = Carbon::now()->addDay();
        } elseif (Carbon::now()->diffInDays($this->article->date) > 0) {
            $nextCrawl = Carbon::now()->addHours(6);
        } else {
            $nextCrawl = Carbon::now()->addHours(2);
        }

        $crawl = $this->article->crawl;
        $crawl->next_crawl = $nextCrawl;

        $this->article->crawl()->save($crawl);
    }
}
