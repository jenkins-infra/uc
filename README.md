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
