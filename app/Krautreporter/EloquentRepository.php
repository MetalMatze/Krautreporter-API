<?php
namespace App\Krautreporter;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\ModelNotFoundException;

abstract class EloquentRepository implements Repository
{
    protected $model;

    public function __construct(Model $model)
    {
        $this->model = $model;
    }

    /**
     * @param array $fields
     * @return \Illuminate\Database\Eloquent\Collection|static[]
     */
    public function all(array $fields = ['*'])
    {
        return $this->model->all($fields);
    }

    /**
     * @param $id
     * @return \Illuminate\Support\Collection|null|static
     * @throws ModelNotFoundException
     */
    public function find($id)
    {
        return $this->model->findOrFail($id);
    }


    /**
     * @param Model $model
     * @return bool
     */
    public function save(Model $model)
    {
        return $model->save();
    }

    /**
     * @param Model $model
     * @return bool|null
     * @throws \Exception
     */
    public function delete(Model $model)
    {
        return $model->delete();
    }
}
