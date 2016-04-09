ghmirror
========

ghmirror helps you to keep copies of your GitHub repositories. It works in two ways:

  * First it's a webhook, which will be called on each push to one of your repository. When it's called it will update its database and local copy.
  * Second, it will be regularly poll GitHub for the list of repositories and update its database and local copies

It's always a good idea to have backups, and `ghmirror` is an ideal solution for backing up your GitHub repositories.

How to install it
-----------------

If you have go installed:

    go get github.com/vrischmann/ghmirror

If not check out the [releases](https://github.com/vrischmann/ghmirror/releases).

How to run it
-------------

First make sure you have the git binary accessible from your PATH: ghmirror depends on it.

You need to set these environment variables one way or another:

  * LISTEN\_ADDRESS               the listen address
  * SECRET                        the secret used by GitHub for the Webhook
  * PERSONAL\_ACCESS\_TOKEN       the token used to authenticate to the GitHub API
  * REPOSITORIES\_PATH            the path where ghmirror will clone the repositories
  * POLL\_FREQUENCY               the frequency at which to poll the repositories list (written as 60s, 1m, 1h, etc)
  * WEBHOOK\_ENDPOINT             the webhook endpoint URL to use when creating a webhook
  * POSTGRES\_HOST                the PostgreSQL hostname
  * POSTGRES\_PORT                the PostgreSQL port
  * POSTGRES\_USER                the PostgreSQL user
  * POSTGRES\_DBNAME              the PostgreSQL database name
  * POSTGRES\_PASSWORD            the PostgreSQL password
  * POSTGRES\_SSLMODE             the PostgreSQL SSL mode (see [here](https://godoc.org/github.com/lib/pq) for valid values)

Your PostgreSQL database needs to have the table defined [here](https://github.com/vrischmann/ghmirror/blob/master/schema.sql). It's up to you to create them one way or another.

The two tables `owner_blacklist` and `repository_blacklist` are used to control which repositories to backup. For example, if you're part of an organization, you may not want to backup their repositories.

Also, right now the cloning of private repositories is not working.

Development
-----------

You need [Go](https://golang.org) 1.6+, [Godep](https://github.com/tools/godep) and [PostgreSQL](http://www.postgresql.org/).

If you don't need to modify one of the dependency, you don't need to do anything: just start coding. Go 1.6+ will always build with the vendored dependencies first.

There's a [Vagrant](https://www.vagrantup.com/) file which will setup a Debian VM with PostgreSQL and create the database with the tables from `schema.sql`.

License
-------

ghmirror is MIT licensed. See the LICENSE file for more details.
