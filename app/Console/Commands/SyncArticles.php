<?php namespace App\Console\Commands;

use App\Article;
use App\Author;
use App\Crawl;
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
     * Array of all parsed articles to later persist in the database
     *
     * @var array
     */
    private $articles;

    /**
     * Create a new command instance.
     *
     * @return void
     */
    public function __construct(Client $client)
    {
        parent::__construct();
        $this->client = $client;
        $this->articles = [];
    }

    /**
     * Execute the console command.
     *
     * @return mixed
     */
    public function fire()
    {
        $this->comment('Begin syncing articles');
        $this->sync_articles();
        $this->save_articles();
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
            $this->parseArticle($node, $index);
        });

        if($nodes->count() > 0)
        {
            $this->comment(count($this->articles));

            $lastArticle = $this->articles[count($this->articles) - 1];
            $url = sprintf('https://krautreporter.de/articles%s/load_more_navigation_items', $lastArticle['url']);
            $this->sync_articles($url);
        }
        else
        {
            $this->info(sprintf('Synced %d articles.', count($this->articles)));
        }
    }

    public function parseArticle(Crawler $node)
    {
        $article = [];

        $a = $node->filter('a');

        $article['url'] = $a->attr('href');

        preg_match('/\/(\d*)/', $article['url'], $matches);
        if(count($matches) >= 2)
        {
            $article['id'] = (int) $matches[1];
        }
        else
        {
            throw new RuntimeException('Failed to parse id from ' . $article['url']);
        }
        $article['author'] = $a->filter('.meta')->text();
        $article['title'] = $a->filter('.item__title')->text();

        if(preg_match('/^Morgenpost:/', $article['title']))
        {
            $article['morgenpost'] = true;
        }
        else
        {
            $article['morgenpost'] = false;
        }

        array_push($this->articles, $article);
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

    private function save_articles()
    {
        $articles = array_reverse($this->articles);

        foreach($articles as $index => $article) {
            $articleModel = Article::firstOrNew(['id' => $article['id']]);

            $articleModel->order = $index;
            $articleModel->title = $article['title'];
            $articleModel->url = $article['url'];
            $articleModel->morgenpost = $article['morgenpost'];

            $author = Author::where('name', '=', $article['author'])->first();

            if($author == null)
            {
                $this->comment('No author found for article ' . $article['url']);
                continue;
            }

            $articleModel->author()->associate($author);
            $articleModel->save();

            if($articleModel->crawl == null)
            {
                $articleModel->crawl()->save(new Crawl());
            }
        }
    }

}
