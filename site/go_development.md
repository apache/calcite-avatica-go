---
layout: docs
title: Go Client Development
sidebar_title: Go Client Development
permalink: /docs/go_development.html
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

## Issues

To file issues, please use the [Calcite JIRA](https://issues.apache.org/jira/projects/CALCITE/issues) and select `avatica-go`
as the component.

## Updating protobuf definitions

To update the procotol buffer definitions, update `AVATICA_VER` in `gen-protobuf.bat` and `gen-protobuf.sh` to match
the version you want to generate protobufs for and then run the appropriate script for your platform.

## Testing

The test suite takes around 4 minutes to run if you run both the Avatica HSQLDB and Apache Phoenix tests.

### Easy way
1. Install [docker](https://docs.docker.com/install/) and [docker-compose](https://docs.docker.com/compose/install/).

2. From the root of the repository, run `docker-compose up`.

### Manual set up
1. Install [Go](https://golang.org/doc/install).

2. Install [dep](https://github.com/golang/dep): `go get -u github.com/golang/dep/cmd/dep`

3. Install dependencies by running `dep ensure -v` from the root of the repository.

4. The test suite requires access to an instance of Avatica running HSQLDB and an instance of Apache Phoenix running the
Phoenix Query Server.

You should then set the `HSQLDB_HOST` and `PHOENIX_HOST` environment variables. For example:
~~~~~~
HSQLDB_HOST: http://hsqldb:8765
PHOENIX_HOST: http://phoenix:8765
~~~~~~

5. To select the test suite, export `AVATICA_FLAVOR=HSQLDB` for Avatica HSQLDB or `AVATICA_FLAVOR=PHOENIX` for Phoenix.

6. Then run `go test -v ./...` from the root of the repository to execute the test suite.

## Releasing
If you have not set up a GPG signing key, set one up by following these [instructions](https://www.apache.org/dev/openpgp.html#generate-key).

From the root of the repository, run `./make-release-artifacts.sh`.

You will be asked to select the tag to build release artifacts for. The latest tag is automatically selected if no tag is selected.

The release artifacts will be placed in the `dist/` folder.