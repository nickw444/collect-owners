Collect OWNERS files and output in a format compatible with Github CODEOWNERS file.


## Usage:
```
usage: Collect Owners [<flags>] <repo>

Walk A Repo and compile a Github CODEOWNERS file

Flags:
  --help                       Show context-sensitive help (also try --help-long and --help-man).
  --contributors=CONTRIBUTORS  Path to contributors file to add to the users DB
  --exclude=EXCLUDE ...        Owners file path exclude patterns
  --add-unresolved             Add ownerships that do not have entries in the users DB as their raw entries within the OWNERS files

Args:
  <repo>  Path to repository
```


## Demo
```
./collect-owners ./demo --add-unresolved
*                       @nickw444/rootowner
component1/*            @octocat
component1/component2/* @nickw444
component2/*            @nickw444
compontent3/*           @octocat
```
