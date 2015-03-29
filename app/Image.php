<?php namespace App;

use Illuminate\Database\Eloquent\Model;

class Image extends Model {

    protected $table = 'images';

    /**
     * No need for timestamps on an image
     *
     * @var bool
     */
    public $timestamps = false;

    public function imageable()
    {
        return $this->morphTo();
    }

}
