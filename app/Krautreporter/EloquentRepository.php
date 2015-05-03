<?php
namespace App\Krautreporter;

use Illuminate\Database\Eloquent\Model;

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
     */
    public function find($id)
    {
        return $this->model->find($id);
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
