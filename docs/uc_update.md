## uc update

uc update --path <path>

### Synopsis

To update all plugins against the latest version of Jenkins:

    uc update --path <path>

To update all plugins against a specific version of Jenkins:

    uc update --path <path> --jenkins-version <version>


```
uc update [flags]
```

### Options

```
      --determine-version-from-dockerfile   Attempt to determine the Jenkins version from a Dockerfile
  -u, --display-updates                     Write updates to stdout
      --dockerfile-path string              Path to the Dockerfile (default "Dockerfile")
  -d, --include-dependencies                Add any additional dependencies to the output
  -j, --jenkins-version string              The version of Jenkins to query against
  -p, --path string                         Path to the plugins.txt file (default "plugins.txt")
  -w, --write                               Update the file rather than display to stdout
```

### Options inherited from parent commands

```
      --help   Show help for command
```

### SEE ALSO

* [uc](uc.md)	 - Update Centre CLI

