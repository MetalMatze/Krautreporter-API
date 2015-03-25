<?php namespace App;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\SoftDeletes;

class Article extends Model {

    use SoftDeletes;

    protected $table = 'articles';

    protected $dates = ['deleted_at'];

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
