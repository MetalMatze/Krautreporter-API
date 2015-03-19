<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\Schema;

class CreateCrawlsTable extends Migration {

    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('crawls', function($table) {
            $table->increments('id');
            $table->integer('crawlable_id')->unsinged();
            $table->string('crawlable_type');
            $table->timestamp('next_crawl');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::drop('crawls');
    }

}
