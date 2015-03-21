@extends('layouts.master')

@section('content')
    @include('crawls._table', ['crawls' => $nextAuthorCrawls])
    @include('crawls._table', ['crawls' => $nextArticleCrawls])
@stop
