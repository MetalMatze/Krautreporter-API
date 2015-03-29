<?php
namespace App\Http\Transformers;

use App\Article;
use League\Fractal\TransformerAbstract;

class ArticleTransformer extends TransformerAbstract {

    protected $availableIncludes = [
        'author',
        'images'
    ];

    protected $defaultIncludes = [
        'images'
    ];

    public function transform(Article $article)
    {
        return [
            'id' => $article->id,
            'order' => $article->order,
            'title' => $article->title,
            'headline' => $article->headline,
            'date' => $article->date,
            'morgenpost' => (bool) $article->morgenpost,
            'url' => $article->url,
            'excerpt' => $article->excerpt,
            'content' => $article->content,
            'author_id' => $article->author_id,
            'created_at' => $article->created_at->format(\DateTime::ISO8601),
            'updated_at' => $article->updated_at->format(\DateTime::ISO8601),
        ];
    }

    public function includeAuthor(Article $article)
    {
        $author = $article->author;

        return $this->item($author, new AuthorTransformer());
    }

    public function includeImages(Article $article)
    {
        $images = $article->images;

        return $this->collection($images, new ImageTransformer());
    }

}
