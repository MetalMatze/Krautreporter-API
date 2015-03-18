<?php

use Illuminate\Database\Migrations\Migration;

class CreateAuthorsTable extends Migration {

    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('authors', function($table)
        {
            $table->integer('id')->unsigend();
            $table->string('name');
            $table->string('title');
            $table->string('url');
            $table->string('image');
            $table->text('biography')->nullable();
            $table->text('media-links')->nullable();
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
        Schema::drop('authors');
    }

}
