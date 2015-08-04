@servers(['production' => 'metalmatze.de'])

@setup
    $repo = 'MetalMatze/Krautreporter-API';

    if ( ! isset($dir) )
    {
        throw new Exception('--dir must be specified');
    }

    $branch      = isset($branch) ? $branch : 'master';
    $release_dir = $dir . '/releases';
    $current_dir = $dir . '/current';
    $release     = date('YmdHis');
@endsetup

@macro('deploy', ['on' => 'production'])
    fetch
    composer
    symlinks
    permissions
    migrate
    clean_old_releases
@endmacro

@task('fetch')
    [ -d {{ $release_dir }} ] || mkdir {{ $release_dir }};
    cd {{ $release_dir }};

    # Make the release dir
    mkdir {{ $release }};

    echo 'Cloning repo';
    git clone --depth=1 git@github.com:{{ $repo }}.git /tmp/{{ $release }};

    # Move the release files
    echo 'Moving release files';
    cd /tmp/{{ $release }}/;
    mv * .[^.]* {{ $release_dir }}/{{ $release }}/;

    # Purge temporary files
    echo 'Purging temporary files';
    rm -rf /tmp/{{ $release }};
@endtask

@task('composer')
    echo 'Installing composer dependencies';
    cd {{ $release_dir }}/{{ $release }};
    /usr/local/bin/composer install --prefer-dist --no-scripts;
    php artisan clear-compiled --env=production;
    php artisan optimize --env=production;
@endtask

@task('symlinks')
    echo 'Updating symlinks';

    # Remove the storage directory and replace with persistent data
    echo 'Linking storage directory';
    rm -rf {{ $release_dir }}/{{ $release }}/storage;
    cd {{ $release_dir }}/{{ $release }};
    ln -nfs ../../storage storage;

    # Import the environment config
    echo 'Linking .env file';
    cd {{ $release_dir }}/{{ $release }};
    ln -nfs ../../.env .env;

    # Symlink the latest release to the current directory
    echo 'Linking current release';
    ln -nfs {{ $release_dir}}/{{$release}} {{ $current_dir }};
@endtask

@task('permissions')
    cd {{ $release_dir }}/{{ $release }};
    echo 'Updating directory permissions';
    find . -type d -exec chmod 775 {} \;
    echo 'Updating file permissions';
    find . -type f -exec chmod 664 {} \;
@endtask

@task('migrate')
    echo 'Running migrations';
    cd {{ $release_dir}}/{{ $release }};
    php artisan migrate --force --env=production;
@endtask

@task('clean_old_releases')
    echo 'Purging old releases';
    # This will list our releases by modification time and delete all but the 5 most recent.
    ls -dt {{ $release_dir }}/* | tail -n +6 | xargs -d '\n' rm -rf;
@endtask
