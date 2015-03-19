<?php
namespace App\Http\Transformers;

use App\Author;
use League\Fractal\TransformerAbstract;

class AuthorTransformer extends TransformerAbstract {

    public function transform(Author $author)
    {
        return [
            'id' => (int) $author->id,
            'name' => $author->name,
            'title' => $author->title,
            'url' => $author->url,
            'image' => $author->image,
            'biography' => $author->biography,
            'socialmedia' => $author->socialmedia,
        ];
    }

}
