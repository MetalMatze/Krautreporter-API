<?php

use Illuminate\Database\Migrations\Migration;

class AddOrderFieldToAuthors extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('authors', function ($table) {
            $table->integer('order')->after('id');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('authors', function ($table) {
            $table->dropColumn('order');
        });
    }
}
