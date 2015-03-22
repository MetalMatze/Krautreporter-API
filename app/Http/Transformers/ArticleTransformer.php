<?php
namespace App\Http\Transformers;

use App\Article;
use League\Fractal\TransformerAbstract;

class ArticleTransformer extends TransformerAbstract {

    public function transform(Article $article)
    {
        return [
            'id' => $article->id,
            'title' => $article->title,
            'headline' => $article->headline,
            'date' => $article->date,
            'morgenpost' => (bool) $article->morgenpost,
            'url' => $article->url,
            'image' => $article->image,
            'excerpt' => $article->excerpt,
            'content' => $article->content,
            'author' => $article->author_id,
            'created_at' => $article->created_at->format(\DateTime::ISO8601),
            'updated_at' => $article->updated_at->format(\DateTime::ISO8601),
        ];
    }

}
