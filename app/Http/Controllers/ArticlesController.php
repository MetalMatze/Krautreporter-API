<?php namespace App\Http\Controllers;

use App\Article;
use App\Http\Requests;
use App\Http\Transformers\ArticleTransformer;
use App\Krautreporter\Articles\ArticleRepository;
use ErrorException;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Support\Facades\Request;
use League\Fractal\Manager;
use League\Fractal\Resource\Collection;
use League\Fractal\Resource\Item;

class ArticlesController extends Controller
{
    /**
     * @var ArticleRepository
     */
    protected $repository;

    /**
     * @var Manager
     */
    protected $fractal;

    /**
     * @var ArticleTransformer
     */
    protected $articleTransformer;

    public function __construct(ArticleRepository $repository, Manager $fractal, ArticleTransformer $articleTransformer)
    {
        $this->repository = $repository;
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
        if (Request::has('olderthan')) {
            $id = (int)Request::get('olderthan');

            try {
                $article = $this->repository->find($id);
            } catch (ModelNotFoundException $e) {
                abort(404);
            }

            $articles = $this->repository->paginateOlderThan($article);
        } else {
            $articles = $this->repository->paginate();
        }

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
        try {
            $article = $this->repository->find($id);
        } catch (ModelNotFoundException $e) {
            abort(404, 'Article not found.');
        }

        $resource = new Item($article, $this->articleTransformer);

        return $this->fractal->createData($resource)->toArray();
    }

}
