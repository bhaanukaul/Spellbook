# Spellbook

## Introduction

Spellbook is a CLI utility to store references for shell commands (aka spells). This is a CRU (Create, Read, Update) app. No "D" since spells cannot be removed from a fantasy realm (afaik).

## Basic Usage

### Initialize Spellbook

Intialize the Spellbook with the default values. 
You can specify an alternative filename for the SQLite DB located at $HOME/.config/Spellbook/Spellbook.db. 
This will also create a config file in the same directory.

```shell
$ Spellbook init [--name]
```

### Create a Spell

```shell
$ Spellbook add --language <language> --content "<spell content>" --description "<description of the spell" --tags tag1,tag2,tag3
```

### Searching for Spells

#### Get all spells in your Spellbook

```shell
$ Spellbook find
ID  Description                   Contents                                Language  Tags
2   list comprehension            x = [row.attribute for row in table]    python    python,lists,comprehension
3   find process bound to a port  netstat -tulpn | egrep "<port number>"  bash      shell,networking
```

#### Get a spell based on ID

```shell
$ Spellbook find id <id>
ID  Description         Contents                              Language  Tags
2   list comprehension  x = [row.attribute for row in table]  python    python,lists,comprehension
```

#### Find spells based on tags

```shell
$ Spellbook find tag <tag>
ID  Description                   Contents                                Language  Tags
3   find process bound to a port  netstat -tulpn | egrep "<port number>"  bash      shell,networking
```

## Advanced Usage

How advanced do you really think this can get?

## Testing

It works on my machine ¯\\\_(ツ)_/¯
