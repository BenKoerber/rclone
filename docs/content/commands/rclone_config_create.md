---
date: 2020-05-27T16:09:49+01:00
title: "rclone config create"
description: "Create a new remote with name, type and options."
slug: rclone_config_create
url: /commands/rclone_config_create/
# autogenerated - DO NOT EDIT, instead edit the source code in cmd/config/create/ and as part of making a release run "make commanddocs"
---
# rclone config create

Create a new remote with name, type and options.

## Synopsis


Create a new remote of `name` with `type` and options.  The options
should be passed in pairs of `key` `value`.

For example to make a swift remote of name myremote using auto config
you would do:

    rclone config create myremote swift env_auth true

Note that if the config process would normally ask a question the
default is taken.  Each time that happens rclone will print a message
saying how to affect the value taken.

If any of the parameters passed is a password field, then rclone will
automatically obscure them if they aren't already obscured before
putting them in the config file.

**NB** If the password parameter is 22 characters or longer and
consists only of base64 characters then rclone can get confused about
whether the password is already obscured or not and put unobscured
passwords into the config file. If you want to be 100% certain that
the passwords get obscured then use the "--obscure" flag, or if you
are 100% certain you are already passing obscured passwords then use
"--no-obscure".  You can also set osbscured passwords using the
"rclone config password" command.

So for example if you wanted to configure a Google Drive remote but
using remote authorization you would do this:

    rclone config create mydrive drive config_is_local false


```
rclone config create `name` `type` [`key` `value`]* [flags]
```

## Options

```
  -h, --help         help for create
      --no-obscure   Force any passwords not to be obscured.
      --obscure      Force any passwords to be obscured.
```

See the [global flags page](/flags/) for global options not listed here.

## SEE ALSO

* [rclone config](/commands/rclone_config/)	 - Enter an interactive configuration session.

