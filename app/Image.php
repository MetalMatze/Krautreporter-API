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

    protected $fillable = [
        'imageable_id',
        'imageable_type',
        'width',
        'src'
    ];

    public function imageable()
    {
        return $this->morphTo();
    }

}
