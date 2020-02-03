# gcr-gc

Keep your Google Container Registry clean :) ... **if you have gcloud SDK installed** ;)

Inspired by https://gist.github.com/ahmetb/7ce6d741bd5baa194a3fac6b1fec8bb7

Install:

    go get github.com/tanelpuhu/gcr-gc

Usage:

    $ gcr-gc -h
    Usage of gcr-gc:
    -r value
            repositories to go over
    -t value
            tags to skip (by default not 'latest')

Example:

    $ gcr-gc -r eu.gcr.io/myrepository
    image eu.gcr.io/myrepository/some-image
    image eu.gcr.io/myrepository/another
     - deleting eu.gcr.io/myrepository/another@sha256:1344b4ef617666ac644d0af5eefb3c1e... created at 2019-11-29 22:55:19+02:00...
    image eu.gcr.io/myrepository/yetone
    ...
