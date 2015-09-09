<?php

use Illuminate\Database\Migrations\Migration;

class CreateArticlesTable extends Migration
{

    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('articles', function ($table) {
            $table->increments('id');
            $table->integer('order')->unsigend();
            $table->string('title');
            $table->string('headline');
            $table->date('date');
            $table->boolean('morgenpost');
            $table->string('url');
            $table->string('image')->nullable();
            $table->text('excerpt')->nullable();
            $table->text('content')->nullable();
            $table->integer('author_id')->unsigned();
            $table->timestamps();
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
