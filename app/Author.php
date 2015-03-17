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
        'homepage',
        'twitter'
    ];

    protected $visible = [
        'id',
        'name',
        'title',
        'url',
        'image',
        'biography',
        'homepage',
        'twitter'
    ];

}
