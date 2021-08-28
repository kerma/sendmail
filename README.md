# I just want to send an email

Sendmail is a a tiny cross-platform CLI tool to send an email via SMTP relay without the need for a configured MTA.

## Usage

    $ sendmail -from charlie@angels.com \
       -to natalie@angels.com,alex@angels.com,dylan@angels.com \
       -cc john@angels.com \
       -attach target.mov \
       -subject "Good morning, angels" < assignment.txt 

## Install

    go install github.com/kerma/sendmail/cmd/sendmail 

## Configuration

In `$XDG_CONFIG_HOME/sendmail/config.json` (or `~/.config/sendmail/config.json`)

    {
        "server": "smtp.angels.com",
        "port": 465,
        "user": "server@angels.com",
        "password": "secret"
    }

Config location can be set via `-conf` flag. Config can also be set via flags: `-server`, `-port`, `-user`, `-password`. Flags will override the ones set in config. If `-user` is left empty no auth will be used. 
