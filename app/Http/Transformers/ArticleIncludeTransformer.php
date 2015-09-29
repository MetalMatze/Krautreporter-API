<?php

namespace App\Http\Transformers;

use App\Article;
use League\Fractal\TransformerAbstract;

class ArticleIncludeTransformer extends TransformerAbstract
{
    public function transform(Article $article)
    {
        return [
            'id' => $article->id,
            'order' => $article->order,
            'title' => $article->title,
            'headline' => $article->headline,
            'date' => $article->date->format(\DateTime::ISO8601),
            'morgenpost' => (bool)$article->morgenpost,
            'preview' => (bool)$article->preview,
            'url' => getenv("URL_KRAUTREPORTER") . $article->url,
            'excerpt' => $article->excerpt,
            'author_id' => $article->author_id,
            'created_at' => $article->created_at->format(\DateTime::ISO8601),
            'updated_at' => $article->updated_at->format(\DateTime::ISO8601),
        ];
    }
}
