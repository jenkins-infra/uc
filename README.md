[![Go Report Card](https://goreportcard.com/badge/github.com/garethjevans/uc)](https://goreportcard.com/report/github.com/garethjevans/uc)
[![Downloads](https://img.shields.io/github/downloads/garethjevans/uc/total.svg)]()

# uc

a small CLI that can be used to update a plugins.txt file for installing plugins into a jenkins docker image.

## To Install

```
brew tap garethjevans/tap
brew install uc
```

This can be used a docker container with the following:

```
docker run -it garethjevans/uc
```

## Usage

### Update Plugins to the latest version

Update plugins.txt to the latest plugin versions, no changes will be made to the file, updates will be pushed to `stdout`.

```
uc update --path plugins.txt
```

Update plugins.txt to the latest plugin versions for the specified jenkins version, no changes will be made to the file, updates will be pushed to `stdout`.

```
uc update --path plugins.txt --jenkins-version 2.263.2
```

Update plugins.txt to the latest plugin versions and automatically update the file

```
uc update --path plugins.txt -w
```

### Only apply security updates

```
uc update --path plugins.txt --security-updates
```

### Determine the Jenkins version from a Dockerfile

`uc` will attempt to determine the Jenkins version from the parent of the Docker image:

```
uc update --path plugins.txt --determine-version-from-dockerfile --dockerfile-path /path/to/Dockerfile
```

### Check for security vulnerabilities

```
uc check --path plugins.txt
```

This will display all known security vulnerabilities for the plugin versions listed:

```
+----------------------+--------------+---------------------------------------------------------------+
| PLUGIN               | ISSUE        | URL                                                           |
+----------------------+--------------+---------------------------------------------------------------+
| github-branch-source | SECURITY-806 | https://jenkins.io/security/advisory/2018-06-04/#SECURITY-806 |
+----------------------+--------------+---------------------------------------------------------------+
```

## Documentation

More indepth documentation can be found [here](./docs/uc.md)

## Development

To build the application:

```
make build
```

To test:

```
make test
```
