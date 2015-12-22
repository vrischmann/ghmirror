ghmirror
========

ghmirror is a webhook server for your GitHub repositories. Its purpose is to clone/pull your GitHub repositories when you push to them so you can always have a local copy.

We all love and use GitHub, there's also not much chance it'll close up or whatever, but better be safe than sorry.

How to install it
-----------------

If you have go installed:

    go get github.com/vrischmann/ghmirror

If not check out the [releases](https://github.com/vrischmann/ghmirror/releases).

How to run it
-------------

First make sure you have the git binary accessible from your PATH: ghmirror depends on it.

You need to set these environment variables one way or another:

  * PORT                          the listen address port
  * SECRET                        the secret used by GitHub for the Webhook
  * PERSONAL\_ACCESS\_TOKEN       the token used to authenticate to the GitHub API
  * POLL\_FREQUENCY               the frequency at which to poll the repositories list (written as 60s, 1m, 1h, etc)
  * WEBHOOK\_ENDPOINT             the webhook endpoint URL to use when creating a webhook
  * WEBHOOK\_VALID\_OWNER\_LOGINS the list of valid owners when creating a webhook. Put your username in here.
  * REPOSITORIES\_PATH            the path where ghmirror will clone the repositories
  * DATABASE\_PATH                the path of the ghmirror database file

ghmirror will regularly poll the GitHub API for the list of repositories, iterate through each one of them and do the following:

  * if it's already in the database, it assumes the webhook has already been installed
    * triggers a git pull
  * otherwise, it checks that the webhook exists
    * if it does not, it creates it
  * then it adds the repository in the database
    * triggers a git clone

Right now there is no way to control which repositories to add the webhook to other than the basic owner login filter.

If there's a need, I will add more ways to filter in the future.

License
-------

ghmirror is MIT licensed. See the LICENSE file for more details.
