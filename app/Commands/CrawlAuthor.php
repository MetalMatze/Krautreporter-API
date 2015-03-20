<?php namespace App\Commands;

use App\Author;
use Carbon\Carbon;
use Goutte\Client;
use Illuminate\Contracts\Bus\SelfHandling;
use Illuminate\Contracts\Queue\ShouldBeQueued;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;
use Symfony\Component\DomCrawler\Crawler;

class CrawlAuthor extends Command implements SelfHandling, ShouldBeQueued {

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

        $crawler->filter('header.article__header')->each(function(Crawler $node)
        {
            $this->author->biography = trim($node->filter('.author__bio')->text());

            try
            {
                $this->author->socialmedia = trim($node->filter('#author-page--media-links')->html());
            }
            catch(\InvalidArgumentException $e) {}

            $this->author->save();

            $crawl = $this->author->crawl;
            $crawl->next_crawl = Carbon::now()->addDay();

            $this->author->crawl()->save($crawl);
        });

    }

}
