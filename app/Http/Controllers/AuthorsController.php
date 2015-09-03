<?php

namespace App\Http\Controllers;

use App\Http\Requests;
use App\Http\Transformers\AuthorTransformer;
use App\Krautreporter\Authors\AuthorRepository;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use League\Fractal\Manager;
use League\Fractal\Resource\Collection;
use League\Fractal\Resource\Item;

class AuthorsController extends Controller
{
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

        $resource = new Collection($authors, $this->authorTransformer);

        return $this->fractal->createData($resource)->toArray();
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

        $resource = new Item($author, $this->authorTransformer);

        return $this->fractal->createData($resource)->toArray();
    }
}
