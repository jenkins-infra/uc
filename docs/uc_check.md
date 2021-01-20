## uc check

uc check --path <path>

### Synopsis

Validate existing plugin versions against known vulnerabilities:

    uc check --path <path>

To check all plugins against a specific version of Jenkins:

    uc check --path <path> --jenkins-version <version>


```
uc check [flags]
```

### Options

```
      --determine-version-from-dockerfile   Attempt to determine the Jenkins version from a Dockerfile
      --dockerfile-path string              Path to the Dockerfile (default "Dockerfile")
  -j, --jenkins-version string              The version of Jenkins to query against
  -p, --path string                         Path to the plugins.txt file (default "plugins.txt")
```

### Options inherited from parent commands

```
  -v, --debug   Debug Output
      --help    Show help for command
```

### SEE ALSO

* [uc](uc.md)	 - Update Centre CLI

