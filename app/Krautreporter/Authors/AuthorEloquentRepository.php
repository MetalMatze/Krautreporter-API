<?php
namespace App\Krautreporter\Authors;

use App\Author;
use App\Krautreporter\EloquentRepository;

class AuthorEloquentRepository extends EloquentRepository implements AuthorRepository
{
    /**
     * @var Author
     */
    protected $model;

    public function __construct(Author $author)
    {
        $this->model = $author;
    }

    public function all(array $fields = ['*'])
    {
        return $this->model->orderBy('order', 'desc')->get($fields);
    }
}
