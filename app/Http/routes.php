<?php

use GuzzleHttp\Client;
use Symfony\Component\DomCrawler\Crawler;

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
    $client = new Client();

    $response = $client->get('https://krautreporter.de/');
    $responseBodyString = $response->getBody()->getContents();

    $crawler = new Crawler($responseBodyString);

    // $register = $crawler->filter('a#registration-link span');
    // dd(trim($register->text()));

    $authors = $crawler->filter('#author-list-tab li a');

    dd(trim($authors->html()));

    foreach ($authors as $index => $author) {
        dd($author->children());
    }
});
