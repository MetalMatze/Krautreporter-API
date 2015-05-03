<?php
namespace App\Krautreporter;

use Illuminate\Database\Eloquent\Model;

interface Repository
{
    public function all(array $fields = []);

    public function find($id);

    public function save(Model $model);

    public function delete(Model $model);
}
