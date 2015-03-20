<?php namespace App\Console\Commands;

use App\Commands\CrawlAuthor;
use App\Crawl;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Queue;

class SyncJobs extends Command {

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
        $jobs = Crawl::where('next_crawl', '<', DB::raw('NOW()'))
                        ->where('crawlable_type', '=', 'App\Author')
                        ->get();

        foreach($jobs as $job)
        {
            Queue::push(new CrawlAuthor($job->crawlable));
        }

        $this->info(sprintf('Added %d jobs to queue', count($jobs)));
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
