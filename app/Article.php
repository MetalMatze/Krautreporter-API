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

    protected $visible = [
        'id',
        'title',
        'date',
        'morgenpost',
        'url',
        'image',
        'excerpt',
        'content',
        'author_id'
    ];

    public function author()
    {
        return $this->belongsTo('App\Author');
    }

}
