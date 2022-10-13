---
layout: news_item
date: "2022-10-13 08:30:00 +0000"
author: francischuang
version: 5.2.0
categories: [release]
tag: v5-2-0
sha: 2b61f37
component: avatica-go
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

Apache Calcite Avatica Go 5.2.0 is a maintenance release of Avatica Go with some dependency updates, bug fixes and a new
minor feature. This release supports Go 1.18 and 1.19, which are currently the versions supported and maintained by the
Go team.

In this release, the `xinsnake/go-http-digest-auth-client` digest authentication client has been replaced with
`icholy/digest` as the former has been deprecated and is no longer maintained. In addition, the `driver.SessionResetter`
has also been implemented, allowing stale connections to the server to be purged and reset.

The Calcite team recommends users of this package to upgrade to this version, where practical, as the dependencies being
used by this package have also been upgraded.