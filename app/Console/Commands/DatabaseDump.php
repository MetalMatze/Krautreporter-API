<?php

namespace App\Console\Commands;

//use App\Helpers\DatabaseMaintenance;
use Ifsnop\Mysqldump\Mysqldump;
use Illuminate\Console\Command;

class DatabaseDump extends Command
{
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
     * @var Mysqldump
     */
    private $dumper;

    /**
     * Create a new command instance.
     *
     * @return void
     */
    public function __construct()
    {
        parent::__construct();

        $this->dumper = new Mysqldump(env('DB_DATABASE'), env('DB_USERNAME'), env('DB_PASSWORD'));
    }

    /**
     * Execute the console command.
     *
     * @return mixed
     */
    public function fire()
    {
        $this->comment('Creating database');

        try {
            $name = DatabaseMaintenance::getBackupName();
            $path = storage_path("app/$name");

            $this->dumper->start($path);

            $this->comment("$path was dumped.");
        } catch (\Exception $e) {
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
