<?php

namespace App\Http\Controllers;

use App\Crawl;
use App\Http\Requests;
use App\Http\Transformers\CrawlTransformer;
use Dingo\Api\Routing\Helpers;
use Illuminate\Http\Response;
use League\Fractal\Manager;
use League\Fractal\Resource\Collection;

class CrawlsController extends Controller
{
    use Helpers;

    /**
     * @var Manager
     */
    private $fractal;

    /**
     * @var CrawlTransformer
     */
    private $crawlTransformer;

    /**
     * CrawlsController constructor.
     * @param Manager $fractal
     * @param CrawlTransformer $crawlTransformer
     */
    public function __construct(Manager $fractal, CrawlTransformer $crawlTransformer)
    {
        $this->fractal = $fractal;
        $this->crawlTransformer = $crawlTransformer;
    }

    /**
     * Display a listing of the resource.
     *
     * @return array
     */
    public function index()
    {
        $crawls = Crawl::with('crawlable')->get();

        return $this->response()->collection($crawls, $this->crawlTransformer);
    }
}
