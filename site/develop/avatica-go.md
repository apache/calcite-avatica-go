---
layout: page
title: Developing the Avatica Go Client
permalink: /develop/avatica-go.html
---

<!--
{% comment %}
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to you under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
{% endcomment %}
-->

* TOC
{:toc}

## Issues

To file issues, please use the [Calcite JIRA](https://issues.apache.org/jira/projects/CALCITE/issues) and select `avatica-go`
as the component.

## Updating protobuf definitions

To update the procotol buffer definitions, update `AVATICA_VER` in `gen-protobuf.bat` and `gen-protobuf.sh` to match
the version you want to generate protobufs for and then run the appropriate script for your platform.

## Live reload during development

It is possible to reload the code in real-time during development. This executes the test suite every time a `.go` or
`.mod` file is updated. The test suite takes a while to run, so the tests will not complete instantly, but live-reloading
during development allows us to not have to manually execute the test suite on save.

### Set up
1. Install [docker](https://docs.docker.com/install/) and [docker-compose](https://docs.docker.com/compose/install/).

2. From the root of the repository, run `DEV=true docker-compose up --build`.

## Testing

The test suite takes around 4 minutes to run if you run both the Avatica HSQLDB and Apache Phoenix tests.

### Easy way
1. Install [docker](https://docs.docker.com/install/) and [docker-compose](https://docs.docker.com/compose/install/).

2. From the root of the repository, run `docker-compose up --build --abort-on-container-exit`.

### Manual set up
1. Install [Go](https://golang.org/doc/install).

For Go 1.10 and below, install the dependencies (skip these steps if using Go 1.11 and above):
1a. Install [dep](https://github.com/golang/dep): `go get -u github.com/golang/dep/cmd/dep`

1b. Install dependencies by running `dep ensure -v` from the root of the repository.

2. The test suite requires access to an instance of Avatica running HSQLDB and an instance of Apache Phoenix running the
Phoenix Query Server.

You should then set the `HSQLDB_HOST` and `PHOENIX_HOST` environment variables. For example:
{% highlight shell %}
HSQLDB_HOST: http://hsqldb:8765
PHOENIX_HOST: http://phoenix:8765
{% endhighlight %}

3. To select the test suite, export `AVATICA_FLAVOR=HSQLDB` for Avatica HSQLDB or `AVATICA_FLAVOR=PHOENIX` for Phoenix.

4. Then run `go test -v ./...` from the root of the repository to execute the test suite.

## Releasing
If you have not set up a GPG signing key, set one up by following these [instructions](https://www.apache.org/dev/openpgp.html#generate-key).

If this release is a new major version (we are releasing 4.0.0 vs the current version 3.0.0), update the version in the
import path in `go.mod`. The import paths in the various sample code snippets should also be updated.

Since we need to support Go modules, tags must be prefixed with a `v`. For example, tag as `v3.1.0` rather than `3.1.0`.

From the root of the repository, run `./make-release-artifacts.sh`.

You will be asked to select the tag to build release artifacts for. The latest tag is automatically selected if no tag is selected.

The release artifacts will be placed in a folder named for the release within the `dist/` folder.

## Important things to note before uploading a release
The name of the release folder must be in the following format: `apache-calcite-avatica-go-$version`. The version must 
include release candidate identifiers such as `-rc0`, if they are present.

The files inside the release folder must have any release candidate identifiers such as `-rc1` removed, even if the
release is a release candidate. `src` must also be added to the filename.

For example, if we are uploading the `apache-calcite-avatica-go-3.0.0-rc1` folder, the files must be named 
`apache-calcite-acatica-go-3.0.0.-srctar.gz`. Note the inclusion of `src` in the filename.

The tar.gz must be named `apache-calcite-avatica-go-$version-src.tar.gz`. 

There must be a GPG signature for the tar.gz named: `apache-calcite-avatica-go-$version-src.tar.gz.asc`

There must be a SHA256 hash for the tar.gz named: `apache-calcite-avatica-go-$version-src.tar.gz.sha256`

## Uploading release artifacts to dev for voting
`svn` must be installed in order to upload release artifacts.

1. Check out the Calcite dev release subdirectory: `svn co "https://dist.apache.org/repos/dist/dev/calcite/" calcite-dev`.

2. Move the release folder under `dist/` into the `calcite-dev` folder.

3. Add the new release to the svn repository: `svn add apache-calcite-avatica-go-3.0.0-rc0`. Remember to change the folder name to the
correct release in the command.

4. Commit to upload the artifacts: `svn commit -m "apache-calcite-avatica-go-3.0.0-rc0" --username yourapacheusername --force-log`
Note the use of `--force-log` to suppress the svn warning, because the commit message is the same as the name of the directory.

## Promoting a release after voting
`svn` must be installed in order to upload release artifacts.

NOTE: Only official releases that has passed a vote may be uploaded to the release directory.

1. Check out the Calcite release directory: `svn co "https://dist.apache.org/repos/dist/release/calcite/" calcite-release`.

2. Copy the release into the `calcite-release` folder. Remember to check the name of the release's folder to ensure that it is in
the correct format.

3. Add the release to the svn repository: `svn add apache-calcite-avatica-go-3.0.0`. Remember to change the folder name to the
correct release in the command.

4. Commit to upload the artifacts: `svn commit -m "Release apache-calcite-avatica-go-3.0.0" --username yourapacheusername`.