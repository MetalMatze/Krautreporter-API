<?php

namespace App\Console\Commands;

use App\Author;
use App\Crawl;
use App\Image;
use Goutte\Client;
use Illuminate\Console\Command;
use Symfony\Component\DomCrawler\Crawler;

class SyncAuthors extends Command
{
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
    protected $description = 'Sync all authors from krautreporter.de.';

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
        $client = new Client();

        $crawler = $client->request('GET', 'http://krautreporter.de');

        $authors = $crawler->filter('#author-list-tab li');

        $this->info(sprintf('Found %d authors, start parsing and saving to db.', count($authors)));

        $authors->each(function (Crawler $node, $index) use ($authors) {

            $anchor = $node->filter('a');

            $author_url = utf8_decode($anchor->attr('href'));

            preg_match('/\/(\d*)/', $author_url, $matches);
            if (count($matches) >= 2) {
                $author_id = (int)$matches[1];
            }

            $author = Author::firstOrNew(['id' => $author_id]);
            $author->url = $author_url;

            $author->order = (count($authors) - $index) - 1;
            $author->name = $anchor->filter('.author__name')->text();

            try {
                $author->title = $anchor->filter('.item__title')->text();
            } catch (\InvalidArgumentException $e) {
            }

            $image = $anchor->filter('img');
            $imageUrls = $image->attr('srcset');

            preg_match('/(.*) 50w, (.*) 100w/', $imageUrls, $matches);
            if (count($matches) == 3) {
                foreach ($matches as $index => $match) {
                    if ($index == 0) {
                        continue;
                    }

                    $image = Image::firstOrCreate([
                        'imageable_id' => $author->id,
                        'imageable_type' => 'App\Author',
                        'width' => $index == 1 ? 130 : 260
                    ]);

                    $image->src = getenv('URL_KRAUTREPORTER') . $match;

                    $author->images()->save($image);
                }
            }

            $author->save();

            if ($author->crawl == null) {
                $crawl = new Crawl();
                $author->crawl()->save($crawl);
            }
        });

    }
}
