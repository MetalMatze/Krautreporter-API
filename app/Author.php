<?php namespace App;

use Illuminate\Database\Eloquent\Model;

class Author extends Model {

    protected $table = 'authors';

    protected $fillable = [
        'id',
        'name',
        'title',
        'url',
        'image',
        'biography',
        'media-links'
    ];

    protected $visible = [
        'id',
        'name',
        'title',
        'url',
        'image',
        'biography',
        'media-links'
    ];

    public function articles()
    {
        return $this->hasMany('App\Article');
    }

    public function crawl()
    {
        return $this->morphOne('App\Crawl', 'crawlable');
    }

}
