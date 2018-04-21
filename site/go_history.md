---
layout: docs
title: Go Client History
permalink: "/docs/go_history.html"
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
[downloads page]({{ site.baseurl }}/downloads/).

## <a href="https://github.com/apache/calcite-avatica-go/releases/tag/3.0.0">3.0.0</a> / 2018-04-23
{: #v3-0-0}

Apache Calcite Avatica Go 3.0.0 is the first release since the Go driver has been donated
to the Apache Software foundation.
Go 1.7+ is supported, but we recommend using the latest stable version of Go.

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

Please refer to the [Boostport/Avatica](https://github.com/Boostport/avatica) Github repository for previous releases
of the Avatica Go client.