<?php

namespace App\Http\Controllers;

use App\Http\Requests;
use App\Http\Transformers\ArticleTransformer;
use App\Krautreporter\Articles\ArticleRepository;
use Dingo\Api\Routing\Helpers;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Support\Facades\Request;
use League\Fractal\Manager;

class ArticlesController extends Controller
{
    use Helpers;

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
     * @return array
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

        return $this->response()->collection($articles, $this->articleTransformer);
    }

    /**
     * Display the specified resource.
     *
     * @param  int $id
     * @return array
     */
    public function show($id)
    {
        try {
            $article = $this->repository->find($id);
        } catch (ModelNotFoundException $e) {
            abort(404);
        }

        return $this->response()->item($article, $this->articleTransformer);
    }
}
