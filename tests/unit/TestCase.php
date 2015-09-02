<?php
namespace tests\unit;

use Mockery;

class TestCase extends \Illuminate\Foundation\Testing\TestCase
{

    public function setUp()
    {
        parent::setUp();

        Mockery::mock('Eloquent');
    }

    public function tearDown()
    {
        parent::tearDown();

        Mockery::close();
    }

    /**
     * Creates the application.
     *
     * @return \Illuminate\Foundation\Application
     */
    public function createApplication()
    {
        $app = require __DIR__.'/../../bootstrap/app.php';

        $app->make(\Illuminate\Contracts\Console\Kernel::class)->bootstrap();

        return $app;
    }
}
