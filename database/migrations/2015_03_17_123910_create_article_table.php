<?php

use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateArticlesTable extends Migration {

    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('articles', function($table)
        {
            $table->integer('id')->unsigned();
            $table->string('title');
            $table->timestamp('date');
            $table->boolean('morgenpost');
            $table->string('url');
            $table->string('image');
            $table->text('excerpt');
            $table->text('content');
            $table->integer('author_id')->unsigned();
            $table->timestamps();

            $table->primary('id');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::drop('articles');
    }

}
