<?php namespace App\Http\Controllers;

use App\Article;
use App\Http\Requests;
use App\Http\Transformers\ArticleTransformer;
use ErrorException;
use Illuminate\Support\Facades\Request;
use League\Fractal\Manager;
use League\Fractal\Resource\Collection;
use League\Fractal\Resource\Item;

class ArticlesController extends Controller
{

    /**
     * @var Manager
     */
    protected $fractal;

    /**
     * @var ArticleTransformer
     */
    protected $articleTransformer;

    public function __construct(Manager $fractal, ArticleTransformer $articleTransformer)
    {
        $this->fractal = $fractal;
        $this->articleTransformer = $articleTransformer;
    }

    /**
     * Display a listing of the resource.
     *
     * @return Response
     */
    public function index()
    {
        $articles = Article::with('images')->orderBy('order', 'desc');

        if (Request::has('olderthan')) {
            try {
                $article = Article::find((int)Request::get('olderthan'));

                $articles = $articles->where('order', '<', $article->order);
            } catch (ErrorException $e) {
            }
        }

        $articles = $articles->take(20)->get();

        $resource = new Collection($articles, $this->articleTransformer);

        return $this->fractal->createData($resource)->toArray();
    }

    /**
     * Display the specified resource.
     *
     * @param  int $id
     * @return Response
     */
    public function show($id)
    {
        $article = Article::with('images')->find($id);

        if ($article == null) {
            abort(404);
        }

        $resource = new Item($article, $this->articleTransformer);

        return $this->fractal->createData($resource)->toArray();
    }

}
