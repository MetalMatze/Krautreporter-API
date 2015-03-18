<?php namespace App\Console\Commands;

use App\Article;
use App\Author;
use Goutte\Client;
use Illuminate\Console\Command;
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
     * Create a new command instance.
     *
     * @return void
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * Execute the console command.
     *
     * @return mixed
     */
    public function fire()
    {
        $client = new Client();

        $crawler = $client->request('GET', 'http://krautreporter.de');

        $nav = $crawler->filter('#article-list-tab li');

        $nav->each(function(Crawler $node, $index) {
            $a = $node->filter('a');

            $article_url = $a->attr('href');

            preg_match('/\/(\d*)/', $article_url, $matches);
            if(count($matches) >= 2)
            {
                $article_id = (int) $matches[1];
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
            $article->order = $index;
            $article->title = $article_title;
            $article->url = $article_url;
            $article->morgenpost = $article_morgenpost;

            $author = Author::where('name', '=', $article_author)->first();
            $article->author()->associate($author);

            $article->save();
        });
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
