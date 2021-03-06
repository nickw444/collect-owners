Walk a Chromium style OWNERS system and output a Github compatible [CODEOWNERS](https://github.com/blog/2392-introducing-code-owners) file

[![Build Status](https://travis-ci.org/nickw444/collect-owners.svg?branch=master)](https://travis-ci.org/nickw444/collect-owners)

## Usage:
```
usage: Collect Owners [<flags>] <repo>

Walk A Repo and compile a Github CODEOWNERS file

Flags:
  --help                       Show context-sensitive help (also try --help-long and --help-man).
  --version                    Show application version.
  --contributors=CONTRIBUTORS  Path to contributors file to add to the users DB
  --exclude=EXCLUDE ...        Regular expressions of paths to exclude
  --add-unresolved             Add ownerships that do not have entries in the users DB as their raw entries within the OWNERS files

Args:
  <repo>  Path to repository

```

## Download

Precompiled binaries are available from Github Releases, [here](//github.com/nickw444/collect-owners/releases)


## Demo
```
./collect-owners ./collect-owners-demo --add-unresolved
*                  @nickw444
component1/*       @nickw444-collect-owners-demo-1
component1/A/*     @nickw444-collect-owners-demo-2
component1/A/*.js  @nickw444-collect-owners-demo-1
component1/B/*     @nickw444-collect-owners-demo-1
component1/C/*     @nickw444-collect-owners-demo-2
component1/D/D.es6 @nickw444-collect-owners-demo-2
```
