<?php

use Illuminate\Database\Migrations\Migration;

class DropImageFromAuthors extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('authors', function ($table) {
            $table->dropColumn('image');
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
            $table->string('image')->nullable();
        });
    }
}
