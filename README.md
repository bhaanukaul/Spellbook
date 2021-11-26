# Spellbook

## Introduction

Spellbook is a CLI utility to store references for shell commands (aka spells). This is a CRU app. No "D" since spells cannot be removed from a fabtasy realm (afaik).

## Basic Usage

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