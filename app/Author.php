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

}
