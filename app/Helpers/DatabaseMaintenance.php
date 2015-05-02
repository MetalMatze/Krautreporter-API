<?php
namespace App\Helpers;

use Carbon\Carbon;

class DatabaseMaintenance
{
    public static function getBackupName()
    {
        $date = Carbon::now()->toDateString();
        return sprintf("%s-%s.sql", env('DB_DATABASE'), $date);
    }
}
