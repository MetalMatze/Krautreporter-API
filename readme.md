Krautreporter-API
====================
[![Author](http://img.shields.io/badge/author-@MetalMatze-blue.svg?style=flat-square)](https://twitter.com/MetalMatze)
[![Latest Version](https://img.shields.io/github/release/MetalMatze/Krautreporter-API.svg?style=flat-square)](https://github.com/MetalMatze/Krautreporter-API/releases)
[![Software License](https://img.shields.io/badge/license-GPLv2-blue.svg?style=flat-square)](https://github.com/MetalMatze/Krautreporter-API/blob/master/LICENSE)
[![Build Status](https://img.shields.io/travis/MetalMatze/Krautreporter-API/master.svg?style=flat-square)](https://travis-ci.org/MetalMatze/Krautreporter-API)

This is a RESTful API written in Laravel PHP for Krautreporter apps.

One app is an Android app created by [MetalMatze](https://github.com/MetalMatze).  
Check it out on github: [MetalMatze/Krautreporter](https://github.com/MetalMatze/Krautreporter)

The other app is an iOS app create by [mindogo GmbH](https://github.com/mindogo).  
Check it out on github: [mindogo/Krautreporter-iOS](https://github.com/mindogo/Krautreporter-iOS)


## Installation

1. Clone this repository to your local machine.
2. Run `$ composer install`.
3. Configure your environment in `.env`.
4. Migrate the database with `php artisan migrate`.
5. Start a standalone server if needed like `php artisan serve`.
6. Visit `localhost:8000` to see the API at work.

## Usage
To fetch data from [krautreporter.de](https://krautreporter.de) you should run the following commands in order

    php artisan sync:authors
    php artisan sync:articles
    php artisan sync:jobs

This will fetch all necessary meta data for authors and articles and create jobs to start crawling.
You can create a cron job to execute this e.g. every 5 minutes.

In addition to running the commands above every once in a while you need also to start the queue worker.

    php artisan queue:listen

This queue worker will listen for new jobs in the database and execute them one after another.

Great! Your database should be able stay up to date with the newest data.

## License

Krautreporter-API is open-source software licensed under the [GPLv2](LICENSE)
