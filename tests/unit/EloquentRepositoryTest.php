<?php
namespace tests\unit;

use App\Krautreporter\EloquentRepository;
use Mockery;

class EloquentRepositoryTest extends TestCase
{
    private $model;
    private $repository;

    public function setUp()
    {
        parent::setUp();

        $this->model = Mockery::mock('Illuminate\Database\Eloquent\Model');
        $this->repository = new DummyEloquentRepository($this->model);
    }

    public function testAll()
    {
        $this->model->shouldReceive('all')->once()->with(['*'])->andReturn(['data']);

        $returned = $this->repository->all();

        $this->assertSame(['data'], $returned);
    }

    public function testAllWithArgs()
    {
        $this->model->shouldReceive('all')->once()->with(['id', 'name'])->andReturn([1, 'foobar']);

        $returned = $this->repository->all(['id', 'name']);

        $this->assertSame([1, 'foobar'], $returned);
    }

    public function testFind()
    {
        $this->model->shouldReceive('find')->once()->with(123)->andReturn($this->model);

        $returned = $this->repository->find(123);

        $this->assertSame($this->model, $returned);
    }

    public function testSave()
    {
        $this->model->shouldReceive('save')->once()->withNoArgs()->andReturn($this->model);

        $returned = $this->repository->save($this->model);

        $this->assertSame($this->model, $returned);
    }

    public function testDelete()
    {
        $this->model->shouldReceive('delete')->once()->withNoArgs()->andReturn($this->model);

        $returned = $this->repository->delete($this->model);

        $this->assertSame($this->model, $returned);
    }
}

class DummyEloquentRepository extends EloquentRepository
{
}
