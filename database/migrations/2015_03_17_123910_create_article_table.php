<?php

use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateArticleTable extends Migration {

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
            $table->string('image');
            $table->string('link');
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
