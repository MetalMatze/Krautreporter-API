<?php

/*
|--------------------------------------------------------------------------
| Model Factories
|--------------------------------------------------------------------------
|
| Here you may define all of your model factories. Model factories give
| you a convenient way to create models for testing and seeding your
| database. Just tell the factory how a default model should look.
|
*/

use Faker\Generator;

$factory->define(App\Author::class, function (Generator $faker) {
    return [
        'name' => $faker->name,
        'title' => $faker->sentence(3),
        'url' => '/' . $faker->slug(6),
        'biography' => $faker->realText(176),
        'socialmedia' => $faker->text(166),
        'created_at' => $faker->dateTimeThisYear,
        'updated_at' => $faker->dateTimeThisMonth
    ];
});

$factory->define(App\Article::class, function (Generator $faker) {
    return [
        'order' => 0,
        'title' => $faker->sentence(6),
        'headline' => $faker->sentence(4),
        'date' => $faker->date('Y-m-d'),
        'morgenpost' => false,
        'preview' => $faker->boolean(23),
        'url' => '/' . $faker->slug(),
        'excerpt' => $faker->text(217),
        'content' => $faker->realText(18233),
        'created_at' => $faker->dateTimeThisYear,
        'updated_at' => $faker->dateTimeThisMonth
    ];
});
