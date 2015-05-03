<?php
namespace App\Krautreporter\Articles;

use App\Article;
use App\Krautreporter\Repository;

interface ArticleRepository extends Repository
{
    public function paginate(array $fields = ['*'], $limit = 20, $direction = 'desc');

    public function paginateOlderThan(Article $article, array $fields = ['*'], $limit = 20, $direction = 'desc');
}
