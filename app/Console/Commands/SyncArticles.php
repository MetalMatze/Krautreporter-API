<?php namespace App\Console\Commands;

use App\Article;
use App\Author;
use Goutte\Client;
use Illuminate\Console\Command;
use RuntimeException;
use Symfony\Component\DomCrawler\Crawler;

class SyncArticles extends Command {

    /**
     * The console command name.
     *
     * @var string
     */
    protected $name = 'sync:articles';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Sync all articles from krautreporter.de.';

    /**
     * @var Client
     */
    protected $client;

    /**
     * Create a new command instance.
     *
     * @return void
     */
    public function __construct(Client $client)
    {
        parent::__construct();
        $this->client = $client;

        $this->lastUrl = null;
        $this->all_index = 0;
    }

    /**
     * Execute the console command.
     *
     * @return mixed
     */
    public function fire()
    {
        $this->sync_articles();
    }

    private function sync_articles($url = null)
    {
        if($url == null)
        {
            $url = 'https://krautreporter.de';
            $filter = '#article-list-tab li';
        }
        else {
            $filter = 'li';
        }

        $crawler = $this->client->request('GET', $url);

        $nodes = $crawler->filter($filter);
        $nodes->each(function(Crawler $node, $index) {
            $a = $node->filter('a');

            $article_url = $a->attr('href');

            preg_match('/\/(\d*)/', $article_url, $matches);
            if(count($matches) >= 2)
            {
                $article_id = (int) $matches[1];
            }
            else
            {
                throw new RuntimeException();
            }

            $article_author = utf8_decode($a->filter('.meta')->text());
            $article_title = $a->filter('.item__title')->text();

            if (preg_match('/^(Morgenpost:)/', $article_title) == 1)
            {
                $article_morgenpost = true;
            }
            else
            {
                $article_morgenpost = false;
            }

            $article = Article::firstOrNew(['id' => $article_id]);
            $article->order = $this->all_index + $index;
            $article->title = $article_title;
            $article->url = $article_url;
            $article->morgenpost = $article_morgenpost;

            $author = Author::where('name', '=', $article_author)->first();
            if($author != null)
            {
                $article->author()->associate($author);
            }
            else
            {
                $this->error(sprintf('Unable to find author with name: %s', $article_author));
            }

            $article->save();

            $this->lastUrl = $article->url;
        });

        if($nodes->count() > 0)
        {
            $this->all_index += $nodes->count();
            $this->comment($this->lastUrl);
            $url = sprintf('https://krautreporter.de/articles%s/load_more_navigation_items', $this->lastUrl);
            $this->sync_articles($url);
        }
    }

    /**
     * Get the console command arguments.
     *
     * @return array
     */
    protected function getArguments()
    {
        return [];
    }

    /**
     * Get the console command options.
     *
     * @return array
     */
    protected function getOptions()
    {
        return [];
    }

}
