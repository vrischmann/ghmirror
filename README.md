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

  * PORT             the listen address port
  * SECRET           the secret used by GitHub for the Webhook
  * REPOSITORIESPATH the path where ghmirror will clone the repositories
  * DATABASEPATH     the path of the ghmirror database file

Now add your webhook through the GitHub interface and point it to your server. The route you need to use is _/hook_

Note: the first time you add the webhook, GitHub will send a ping event: ghmirror does nothing in that case. Right now ghmirror only handles push events.

Future work
-----------

 * automatically creates webhooks via the GitHub API
   * fetch all repositories regularly
   * create the webhook if not present
 * handle all webhook events maybe ?

License
-------

ghmirror is MIT licensed. See the LICENSE file for more details.
