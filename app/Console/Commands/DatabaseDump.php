<?php namespace App\Console\Commands;

use Carbon\Carbon;
use Ifsnop\Mysqldump\Mysqldump;
use Illuminate\Console\Command;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Input\InputArgument;

class DatabaseDump extends Command {

    /**
     * The console command name.
     *
     * @var string
     */
    protected $name = 'db:dump';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Create a dump of the database.';

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
        $this->comment(sprintf("Creating database dump for %s", env('DB_DATABASE')));

        try
        {
            $date = Carbon::now()->format('Y-m-d');
            $name = sprintf("%s-%s.sql", env('DB_DATABASE'), $date);

            $dump = new Mysqldump(env('DB_DATABASE'), env('DB_USERNAME'), env('DB_PASSWORD'));
            $dump->start(sprintf("storage/app/%s", $name));
        }
        catch (\Exception $e)
        {
            $this->error($e->getMessage());
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
