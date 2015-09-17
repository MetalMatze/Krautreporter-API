<?php

namespace App\Http\Controllers;

use App\Crawl;
use App\Http\Requests;
use Illuminate\Support\Facades\Response;

class StatusController extends Controller
{
    /**
     * Display a listing of the resource.
     *
     * @return Response
     */
    public function index()
    {
        $nextAuthorCrawls = Crawl::where('crawlable_type', '=', 'App\Author')
            ->orderBy('next_crawl', 'asc')
            ->limit(10)
            ->get();

        $nextArticleCrawls = Crawl::where('crawlable_type', '=', 'App\Article')
            ->orderBy('next_crawl', 'asc')
            ->limit(10)
            ->get();

        return view('status')
            ->with('nextAuthorCrawls', $nextAuthorCrawls)
            ->with('nextArticleCrawls', $nextArticleCrawls);
    }
}
