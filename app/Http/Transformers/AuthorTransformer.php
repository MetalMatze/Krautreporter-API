<?php

namespace App\Http\Transformers;

use App\Author;
use League\Fractal\TransformerAbstract;

class AuthorTransformer extends TransformerAbstract
{
    protected $availableIncludes = [
        'images'
    ];

    protected $defaultIncludes = [
        'images'
    ];

    public function transform(Author $author)
    {
        return [
            'id' => (int)$author->id,
            'name' => $author->name,
            'title' => $author->title,
            'url' => getenv("URL_KRAUTREPORTER") . $author->url,
            'biography' => $author->biography,
            'socialmedia' => $author->socialmedia,
            'created_at' => $author->created_at->format(\DateTime::ISO8601),
            'updated_at' => $author->updated_at->format(\DateTime::ISO8601),
        ];
    }

    public function includeImages(Author $author)
    {
        $images = $author->images;

        return $this->collection($images, new ImageTransformer());
    }
}
