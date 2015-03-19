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
            'date' => $article->date,
            'morgenpost' => (bool) $article->morgenpost,
            'url' => $article->url,
            'image' => $article->image,
            'excerpt' => $article->excerpt,
            'contect' => $article->content,
            'author' => $article->author_id
        ];
    }

}
