<?php

namespace App;

use Illuminate\Database\Eloquent\Model;

class Crawl extends Model
{
    protected $table = 'crawls';

    public $timestamps = false;

    public function crawlable()
    {
        return $this->morphTo();
    }
}
