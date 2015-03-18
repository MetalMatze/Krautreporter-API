<?php

use App\Article;
use App\Author;

/*
|--------------------------------------------------------------------------
| Application Routes
|--------------------------------------------------------------------------
|
| Here is where you can register all of the routes for an application.
| It's a breeze. Simply tell Laravel the URIs it should respond to
| and give it the controller to call when that URI is requested.
|
*/

Route::get('/', function() {
    return 'hallo';
});

Route::get('authors', function() {
    return Author::all();
});

Route::get('authors/{id}', function($id) {
    return Author::find($id);
});

Route::get('articles', function() {
    return Article::orderBy('order', 'asc')->get();
});
