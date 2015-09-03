<?php

namespace App\Console;

use Illuminate\Console\Scheduling\Schedule;
use Illuminate\Foundation\Console\Kernel as ConsoleKernel;

class Kernel extends ConsoleKernel
{
    /**
     * The Artisan commands provided by your application.
     *
     * @var array
     */
    protected $commands = [
        \App\Console\Commands\Sync::class,
        \App\Console\Commands\SyncJobs::class,
        \App\Console\Commands\SyncAuthors::class,
        \App\Console\Commands\DatabaseDump::class,
        \App\Console\Commands\SyncArticles::class,
        \App\Console\Commands\DatabaseBackup::class,
    ];

    /**
     * Define the application's command schedule.
     *
     * @param  \Illuminate\Console\Scheduling\Schedule $schedule
     * @return void
     */
    protected function schedule(Schedule $schedule)
    {
        $schedule->command('sync:authors && php artisan sync:articles && php artisan sync:jobs')->everyFiveMinutes();
        $schedule->command('db:dump && php artisan db:backup')->daily();
    }
}
