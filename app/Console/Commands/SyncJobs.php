<?php

namespace App\Console\Commands;

//use App\Commands\CrawlArticle;
//use App\Commands\CrawlAuthor;
use App\Crawl;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Queue;

class SyncJobs extends Command
{
    /**
     * The console command name.
     *
     * @var string
     */
    protected $name = 'sync:jobs';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Command description.';

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
        $authorJobs = Crawl::where('next_crawl', '<', DB::raw('NOW()'))
            ->where('crawlable_type', '=', 'App\Author')
            ->orderBy('next_crawl', 'asc')
            ->orderBy('crawlable_id', 'desc')
            ->get();

        foreach ($authorJobs as $job) {
            Queue::push(new CrawlAuthor($job->crawlable));
        }

        $articlesJobs = Crawl::where('next_crawl', '<', DB::raw('NOW()'))
            ->where('crawlable_type', '=', 'App\Article')
            ->orderBy('next_crawl', 'asc')
            ->orderBy('crawlable_id', 'desc')
            ->get();

        foreach ($articlesJobs as $job) {
            Queue::push(new CrawlArticle($job->crawlable));
        }

        $this->info(sprintf('Added %d jobs to queue', count($authorJobs) + count($articlesJobs)));
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
