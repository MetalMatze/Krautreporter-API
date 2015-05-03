<?php
namespace App\Krautreporter\Articles;

use App\Article;
use App\Krautreporter\EloquentRepository;
use Illuminate\Database\Eloquent\Model;

class ArticleEloquentRepository extends EloquentRepository implements ArticleRepository
{
    /**
     * @var Article
     */
    protected $model;

    public function __construct(Article $model)
    {
        $this->model = $model;
    }

    public function paginate(array $fields = ['*'], $limit = 20, $direction = 'desc')
    {
        return $this->articlePagination($limit, $direction)->get($fields);
    }

    public function paginateOlderThan(Article $article, array $fields = ['*'], $limit = 20, $direction = 'desc')
    {
        return $this->articlePagination($limit, $direction)
            ->where('order', '<', $article->order)
            ->get($fields);
    }

    /**
     * @param $limit
     * @param $direction
     * @return mixed
     */
    private function articlePagination($limit, $direction)
    {
        return $this->model
            ->with('images')
            ->orderBy('order', $direction)
            ->take($limit);
    }
}
