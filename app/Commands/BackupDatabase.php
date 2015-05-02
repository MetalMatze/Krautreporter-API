<?php namespace App\Commands;

use App\Helpers\DatabaseMaintenance;
use Exception;
use Illuminate\Contracts\Bus\SelfHandling;
use Illuminate\Contracts\Filesystem\Factory;
use Illuminate\Support\Facades\Storage;

class BackupDatabase extends Command implements SelfHandling
{
    private $dropboxDisk;
    private $localDisk;

    /**
     * Create a new command instance.
     *
     * @param Factory $storage
     */
    public function __construct()
    {
        $this->localDisk = Storage::disk('local');
        $this->dropboxDisk = Storage::disk('dropbox');
    }

    /**
     * Execute the command.
     *
     * @return void
     */
    public function handle()
    {
        $name = DatabaseMaintenance::getBackupName();

        $backupContent = $this->getBackupContent($name);

        $this->backup($name, $backupContent);
    }

    /**
     * @param $name
     * @return mixed
     * @throws Exception
     */
    private function getBackupContent($name)
    {
        if (!$this->localDisk->exists($name)) {
            throw new Exception("Cannot backup $name");
        }

        $backupContent = $this->localDisk->get($name);

        return $backupContent;
    }

    /**
     * @param $name
     * @param $backupContent
     * @throws Exception
     */
    private function backup($name, $backupContent)
    {
        $successful = $this->dropboxDisk->put($name, $backupContent);

        if (!$successful) {
            throw new Exception('Could not write backup');
        }
    }

}
