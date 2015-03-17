<?php namespace App\Console\Commands;

use App\Author;
use GuzzleHttp\Client;
use Illuminate\Console\Command;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\DomCrawler\Crawler;

class SyncAuthors extends Command {

    /**
     * The console command name.
     *
     * @var string
     */
    protected $name = 'sync:authors';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Synchronize (create, update) all authors from krautreporter website.';

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

        $response = $client->get('https://krautreporter.de/');
        $responseBodyString = $response->getBody()->getContents();

        $crawler = new Crawler($responseBodyString);

        $authors = $crawler->filter('#author-list-tab li');

        $authors->each(function(Crawler $node) {

            $anchor = $node->filter('a');

            $author_url = $anchor->attr('href');

            preg_match('/\/(\d*)/', $author_url, $matches);
            if(count($matches) >= 2)
            {
                $author_id = (int) $matches[1];
            }

            $author = Author::firstOrNew(['id' => $author_id]);

            $anchorText = trim($anchor->text());

            preg_match('/(.*)(\n(.*))?/', $anchorText, $matches);

            if(count($matches) >= 2)
            {
                $author->name = $matches[1];
            }

            if (count($matches) >= 4)
            {
                $author->title = $matches[3];
            }

            $image = $anchor->filter('img');
            $imageUrls = $image->attr('srcset');

            preg_match('/(.*) 50w, (.*) 100w/', $imageUrls, $matches);
            if(count($matches) == 3)
            {
                $author->image = $matches[2];
            }

            $author->save();
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
