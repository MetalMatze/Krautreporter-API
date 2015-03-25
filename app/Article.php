<?php namespace App;

use Illuminate\Database\Eloquent\Model;

class Article extends Model {

    protected $table = 'articles';

    protected $fillable = [
        'id',
        'title',
        'date',
        'morgenpost',
        'url',
        'image',
        'excerpt',
        'content'
    ];

    public function author()
    {
        return $this->belongsTo('App\Author');
    }

    public function crawl()
    {
        return $this->morphOne('App\Crawl', 'crawlable');
    }

}
