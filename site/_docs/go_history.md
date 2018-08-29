---
layout: docs
title: Go Client History
permalink: /docs/go_history.html
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

For a full list of releases, see
<a href="https://github.com/apache/calcite-avatica-go/releases">github</a>.
Downloads are available on the
[downloads page]({{ site.baseurl }}/downloads/avatica-go.html).

## <a href="https://github.com/apache/calcite-avatica-go/releases/tag/v3.1.0">3.1.0</a> / 2018-04-27
{: #v3-1-0}

Apache Calcite Avatica Go 3.1.0 is a minor release of Avatica Go with support for Go modules.
We recommend using the latest stable version of Go.

Go 1.11 along with Go modules support was released at the end of August 2018. Go modules will become the
official package management solution for Go projects. As the Go team currently supports both Go 1.11 and Go 1.10,
the Gopkg.toml and Gopkg.lock files are still available for those using dep for package management. We plan to
remove support for dep when Go 1.12 is released in early 2019, so we encourage users to upgrade to Go 1.11 and use
Go modules where possible.

Features and bug fixes

* [<a href="https://issues.apache.org/jira/browse/CALCITE-2333">CALCITE-2333</a>]
  Stop releasing zip archives
* [<a href="https://issues.apache.org/jira/browse/CALCITE-2367">CALCITE-2367</a>]
  Remove the satori/go.uuid package as it is unmaintained and replace it with hashicorp/go-uuid which is already
  a transitive dependency (Kenneth Shaw)
* [<a href="https://issues.apache.org/jira/browse/CALCITE-2372">CALCITE-2372</a>]
  Test against Apache Phoenix 4.14.0
* [<a href="https://issues.apache.org/jira/browse/CALCITE-2493">CALCITE-2493</a>]
  Update all dependencies to their latest versions
* [<a href="https://issues.apache.org/jira/browse/CALCITE-2335">CALCITE-2335</a>]
  Add support for Go modules (available since Go 1.11) and test against Go 1.11
* [<a href="https://issues.apache.org/jira/browse/CALCITE-2500">CALCITE-2500</a>]
  Test against Avatica 1.12.0 and Apache Phoenix 5.0.0 and regenerate protobuf definitions
* Fix release script

Web site and documentation

* [<a href="https://issues.apache.org/jira/browse/CALCITE-2335">CALCITE-2335</a>]
  Update documentation and release instructions to include support for Go modules.

## <a href="https://github.com/apache/calcite-avatica-go/releases/tag/3.0.0">3.0.0</a> / 2018-04-27
{: #v3-0-0}

Apache Calcite Avatica Go 3.0.0 is the first release since the Go driver has been donated
to the Apache Software foundation.
We recommend using the latest stable version of Go.

Features and bug fixes

* [<a href="https://issues.apache.org/jira/browse/CALCITE-1240">CALCITE-1240</a>]
  Intial import of the original [Boostport/avatica](https://github.com/Boostport/avatica) code-base into the
  [apache/calcite-avatica-go](https://github.com/apache/calcite-avatica-go) repository
* [<a href="https://issues.apache.org/jira/browse/CALCITE-1938">CALCITE-1938</a>]
  Releasing the first release of Calcite Avatica Go under the Apache Software Foundation
* Remove go-cleanhttp dependency
* Support for Avatica HSQLDB backend and move Apache Phoenix support into adapter
* Add bash script to automate releases with checks to alert on files without the Apache license header
* Replace gopher.png test fixture with Calcite logo

Web site and documentation

* [<a href="https://issues.apache.org/jira/browse/CALCITE-1937">CALCITE-1937</a>]
  Set up Calcite Avatica Go website

## Past releases

Prior to release 3.0.0, the Avatica Go client was developed by Boostport.

Please refer to the [Boostport/avatica](https://github.com/Boostport/avatica) Github repository for previous releases
of the Avatica Go client.