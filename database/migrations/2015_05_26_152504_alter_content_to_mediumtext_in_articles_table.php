<?php

use Illuminate\Database\Migrations\Migration;

class AlterContentToMediumtextInArticlesTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        DB::statement('ALTER TABLE articles MODIFY content MEDIUMTEXT');
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        DB::statement('ALTER TABLE articles MODIFY content TEXT');
    }
}
