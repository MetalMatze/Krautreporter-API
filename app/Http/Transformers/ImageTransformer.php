<?php
namespace App\Http\Transformers;

use App\Image;
use League\Fractal\TransformerAbstract;

class ImageTransformer extends TransformerAbstract {


    public function transform(Image $image)
    {
        return [
            'width' => $image->width,
            'src' => $image->src
        ];
    }

}

