<?php namespace App\Console\Commands;

use App\Article;
use App\Author;
use App\Crawl;
use Carbon\Carbon;
use Goutte\Client;
use Illuminate\Console\Command;
use RuntimeException;
use Symfony\Component\DomCrawler\Crawler;

class SyncArticles extends Command
{
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

    /**
     * Execute the console command.
     *
     * @return mixed
     */
    public function fire()
    {
        $this->comment('Begin syncing articles');
        $this->syncArticles();
        $this->saveArticles();
    }

    private function syncArticles($url = null)
    {
        if ($url == null) {
            $url = 'https://krautreporter.de';
            $filter = '#article-list-tab li';
        } else {
            $filter = 'li';
        }

        $crawler = $this->client->request('GET', $url);

        $nodes = $crawler->filter($filter);
        $nodes->each(function (Crawler $node, $index) {
            $this->parseArticle($node, $index);
        });

        if ($nodes->count() > 0) {
            $this->comment(count($this->articles));

            $lastArticle = $this->articles[count($this->articles) - 1];
            $url = sprintf('https://krautreporter.de/articles%s/load_more_navigation_items', $lastArticle['url']);
            $this->syncArticles($url);
        } else {
            $this->info(sprintf('Synced %d articles.', count($this->articles)));
        }
    }

    public function parseArticle(Crawler $node)
    {
        $article = [];

        $a = $node->filter('a');

        $article['url'] = $a->attr('href');
        $article['id'] = $this->parseId($article['url']);
        $article['author'] = $a->filter('.meta')->text();
        $article['title'] = $a->filter('.item__title')->text();
        $article['morgenpost'] = $this->isMorgenpost($article['title']);
        $article['preview'] = $this->isPreview($a);

        array_push($this->articles, $article);
    }

    /**
     * @param $title
     * @return mixed
     */
    private function isMorgenpost($title)
    {
        if (preg_match('/^Morgenpost:/', $title)) {
            return true;
        }

        return false;
    }

    /**
     * @param $url
     * @param $matches
     * @return mixed
     */
    private function parseId($url)
    {
        preg_match('/\/(\d*)/', $url, $matches);
        if (count($matches) >= 2) {
            return (int)$matches[1];
        } else {
            throw new RuntimeException('Failed to parse id from ' . $url);
        }
    }

    private function isPreview(Crawler $a)
    {
        try {
            if ($a->filter('img')->count() == 1) {
                return true;
            }
        } catch (\InvalidArgumentException $e) {
        }

        return false;
    }

    private function saveArticles()
    {
        $articles = array_reverse($this->articles);

        foreach ($articles as $index => $article) {
            $articleModel = Article::withTrashed()->where('id', '=', $article['id'])->first();

            if ($articleModel == null) {
                $articleModel = new Article();
                $articleModel->id = $article['id'];
                $articleModel->date = Carbon::now();
            }

            $articleModel->order = $index;
            $articleModel->title = $article['title'];
            $articleModel->url = $article['url'];
            $articleModel->morgenpost = $article['morgenpost'];
            $articleModel->preview = $article['preview'];

            $author = Author::where('name', '=', $article['author'])->first();

            if ($author == null) {
                $this->comment('No author found for article ' . $article['url']);
                continue;
            }

            $articleModel->author()->associate($author);
            $articleModel->save();

            if ($articleModel->crawl == null) {
                $articleModel->crawl()->save(new Crawl());
            }
        }
    }
}
