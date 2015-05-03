<?php
namespace App\Krautreporter;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\ModelNotFoundException;

interface Repository
{
    /**
     * @param array $fields
     * @return mixed
     */
    public function all(array $fields = []);

    /**
     * @param $id
     * @return Model
     * @throws ModelNotFoundException
     */
    public function find($id);

    /**
     * @param Model $model
     * @return bool
     */
    public function save(Model $model);

    /**
     * @param Model $model
     * @return mixed
     * @throws \Exception
     */
    public function delete(Model $model);
}
