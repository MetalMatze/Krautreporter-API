<?php namespace App\Http\Controllers;

use App\Author;
use App\Http\Requests;
use App\Http\Transformers\AuthorTransformer;
use Illuminate\Http\Response;
use League\Fractal\Manager;
use League\Fractal\Resource\Collection;
use League\Fractal\Resource\Item;

class AuthorsController extends Controller {

    /**
     * @var Manager
     */
    protected $fractal;

    /**
     * @var AuthorTransformer
     */
    protected $authorTransformer;

    /**
     * @param Manager $fractal
     * @param AuthorTransformer $authorTransformer
     */
    function __construct(Manager $fractal, AuthorTransformer $authorTransformer)
    {
        $this->fractal = $fractal;
        $this->authorTransformer = $authorTransformer;
    }

    /**
     * Display a listing of the resource.
     *
     * @return Response
     */
    public function index()
    {
        $authors = Author::all();

        $resource = new Collection($authors, $this->authorTransformer);

        return $this->fractal->createData($resource)->toArray();
    }

    /**
     * Display the specified resource.
     *
     * @param  int  $id
     * @return Response
     */
    public function show($id)
    {
        $author = Author::find($id);

        if($author == null) {
            abort(404);
        }

        $resource = new Item($author, $this->authorTransformer);

        return $this->fractal->createData($resource)->toArray();
    }

}
