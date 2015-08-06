<?php namespace App\Commands;

use App\Author;
use App\Image;
use Carbon\Carbon;
use Goutte\Client;
use Illuminate\Contracts\Bus\SelfHandling;
use Illuminate\Contracts\Queue\ShouldBeQueued;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;
use Symfony\Component\DomCrawler\Crawler;

class CrawlAuthor extends Command implements SelfHandling, ShouldBeQueued
{
    use InteractsWithQueue, SerializesModels;

    /**
     * @var Author
     */
    private $author;

    /**
     * Create a new command instance.
     *
     * @return void
     */
    public function __construct(Author $author)
    {
        $this->author = $author;
    }

    /**
     * Execute the command.
     *
     * @return void
     */
    public function handle()
    {
        $client = new Client();

        $crawler = $client->request('GET', 'https://krautreporter.de' . $this->author->url);

        if ($client->getResponse()->getStatus() == 200) {
            $crawler->filter('header.article__header')->each(function (Crawler $node) {
                $this->author->biography = trim($node->filter('.author__bio')->text());

                if ($this->author->biography == '') {
                    $this->author->biography = null;
                }

                try {
                    $this->author->socialmedia = trim($node->filter('#author-page--media-links')->html());
                } catch (\InvalidArgumentException $e) {
                }

                $this->author->save();

                $imageUrls = $node->filter('h2.author--large img')->attr('srcset');
                preg_match('/(.*) 170w, (.*) 340w/', $imageUrls, $matches);
                if (count($matches) == 3) {
                    foreach ($matches as $index => $match) {
                        if ($index == 0) {
                            continue;
                        }

                        $image = Image::firstOrCreate([
                            'imageable_id' => $this->author->id,
                            'imageable_type' => 'App\Author',
                            'width' => $index == 1 ? 170 : 340
                        ]);

                        $image->src = getenv('URL_KRAUTREPORTER') . $match;

                        $this->author->images()->save($image);
                    }
                }

                $crawl = $this->author->crawl;
                $crawl->next_crawl = Carbon::now()->addDay();

                $this->author->crawl()->save($crawl);
            });
        }
    }
}
