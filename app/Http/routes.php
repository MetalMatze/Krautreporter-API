<?php

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


Route::get('/', function () {
    return 'hallo';
});

Route::get('status', 'StatusController@index');

$api = app('Dingo\Api\Routing\Router');


$api->version('v1', function ($api) {
    $api->get('authors', 'App\Http\Controllers\AuthorsController@index');
    $api->get('authors/{id}', 'App\Http\Controllers\AuthorsController@show');

    $api->get('articles', 'App\Http\Controllers\ArticlesController@index');
    $api->get('articles/{id}', 'App\Http\Controllers\ArticlesController@show');
});
