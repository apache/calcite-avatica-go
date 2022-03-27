---
layout: news_item
date: "2022-03-27 08:30:00 +0000"
author: francischuang
version: 5.1.0
categories: [release]
tag: v5-1-0
sha: c232d7b
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

Apache Calcite Avatica Go 5.1.0 is a maintenance release of Avatica Go with some dependency updates and bug fixes.
This release supports Go 1.17 and 1.18, which are currently the versions supported and maintained by the Go team.

Of particular note is the replacement of the [github.com/golang/protobuf](https://github.com/golang/protobuf) package with
[google.golang.org/protobuf](https://pkg.go.dev/google.golang.org/protobuf) and the 
[golang.org/x/xerrors](https://pkg.go.dev/golang.org/x/xerrors) package with the standard library's built-in 
[errors](https://pkg.go.dev/errors) package as the former packages have all be superseded by their replacements.

The Calcite team recommends users of this package to upgrade to this version, where practical, as the dependencies being
used by this package have also been upgraded.