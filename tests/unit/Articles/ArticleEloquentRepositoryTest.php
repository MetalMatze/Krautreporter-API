<?php
namespace tests\unit\Articles;

use App\Krautreporter\Articles\ArticleEloquentRepository;
use Mockery;
use tests\unit\TestCase;

class ArticleEloquentRepositoryTest extends TestCase
{
    protected $model;
    protected $repository;

    public function setUp()
    {
        parent::setUp();

        $this->model = Mockery::mock('\App\Article');
        $this->repository = new ArticleEloquentRepository($this->model);
    }

    public function testPaginate()
    {
        $this->mockArticlePagination(20, 'desc')
            ->shouldReceive('get')->once()->with(['*'])->andReturn(['foo', 'bar']);

        $returned = $this->repository->paginate();

        $this->assertSame(['foo', 'bar'], $returned);
    }

    public function testPaginateOlderThan()
    {
        $this->model->shouldReceive('getAttribute')->once()->with('order')->andReturn(321);

        $this->mockArticlePagination(20, 'desc')
            ->shouldReceive('where')->once()->with('order', '<', 321)->andReturn($this->model)
            ->shouldReceive('get')->once()->with(['*'])->andReturn(['foo', 'bar']);

        $returned = $this->repository->paginateOlderThan($this->model);

        $this->assertSame(['foo', 'bar'], $returned);
    }

    private function  mockArticlePagination($limit, $direction)
    {
        return $this->model->shouldReceive('with')->with('images')->once()->andReturn($this->model)
            ->shouldReceive('orderBy')->once()->with('order', $direction)->andReturn($this->model)
            ->shouldReceive('take')->once()->with($limit)->andReturn($this->model);
    }

}
