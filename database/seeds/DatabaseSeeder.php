<?php

use App\Article;
use App\Author;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;

class DatabaseSeeder extends Seeder
{
    /**
     * Run the database seeds.
     *
     * @return void
     */
    public function run()
    {
        Model::unguard();

        factory(Author::class)->times(10)->create();

        DB::beginTransaction();
        foreach (range(0, 99) as $index) {
            $article = factory(Article::class)->make();
            $article->author()->associate(Author::all()->random(1));
            $article->order = $index;
            $article->save();
        }
        DB::commit();

        Model::reguard();
    }
}
