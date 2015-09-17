<?php

namespace App\Http\Controllers;

use App\Http\Requests;
use App\Http\Transformers\AuthorTransformer;
use App\Krautreporter\Authors\AuthorRepository;
use Dingo\Api\Routing\Helpers;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use League\Fractal\Manager;

class AuthorsController extends Controller
{
    use Helpers;

    /**
     * @var AuthorRepository
     */
    protected $repository;

    /**
     * @var Manager
     */
    protected $fractal;

    /**
     * @var AuthorTransformer
     */
    protected $authorTransformer;

    /**
     * @param AuthorRepository $repository
     * @param Manager $fractal
     * @param AuthorTransformer $authorTransformer
     */
    public function __construct(AuthorRepository $repository, Manager $fractal, AuthorTransformer $authorTransformer)
    {
        $this->repository = $repository;
        $this->fractal = $fractal;
        $this->authorTransformer = $authorTransformer;
    }

    /**
     * Display a listing of the resource.
     *
     * @return array
     */
    public function index()
    {
        $authors = $this->repository->all();

        return $this->response()->collection($authors, $this->authorTransformer);
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
            $author = $this->repository->find($id);
        } catch (ModelNotFoundException $e) {
            abort(404);
        }

        return $this->response()->item($author, $this->authorTransformer);
    }
}
