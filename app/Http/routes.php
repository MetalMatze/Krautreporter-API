<?php

use App\Author;
use App\Article;
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
    return 'hallo';
});

Route::get('authors', function() {
    return Author::all();
});

Route::get('authors/{id}', function($id) {
    return Author::find($id);
});

Route::get('articles', function() {
    return Article::orderBy('id', 'desc')->get();
});

Route::get('sync', function() {

    $client = new Client();

    $response = $client->get('https://krautreporter.de/');
    $responseBodyString = $response->getBody()->getContents();

    $crawler = new Crawler($responseBodyString);

    $nav = $crawler->filter('#article-list-tab li');

    $nav->each(function(Crawler $node) {
        // dd($node->html());
        $a = $node->filter('a');
        $article_url = $a->attr('href');

        preg_match('/\/(\d*)/', $article_url, $matches);
        if(count($matches) >= 2)
        {
            $article_id = (int) $matches[1];
        }

        var_dump($a->html());

        $article_author = utf8_decode($a->filter('.meta')->text());
        $article_title = utf8_decode($a->filter('.item__title')->text());

        if (preg_match('/^(Morgenpost:)/', $article_title) == 1)
        {
            $article_morgenpost = true;
        }
        else
        {
            $article_morgenpost = false;
        }

        $article = Article::firstOrNew(['id' => $article_id]);
        $article->title = $article_title;
        $article->url = $article_url;
        $article->morgenpost = $article_morgenpost;

        $author = Author::where('name', '=', $article_author)->first();
        $article->author()->associate($author);

        $article->save();
    });
});