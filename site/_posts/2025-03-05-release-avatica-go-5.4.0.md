---
layout: news_item
date: "2025-03-05 08:30:00 +0000"
author: francischuang
version: 5.4.0
categories: [release]
tag: v5-4-0
sha: f0d392b
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

Apache Calcite Avatica Go 5.4.0 includes dependency updates and various minor improvements.

The Avatica protobuf messages have been recompiled to use the new [opaque API](https://go.dev/blog/protobuf-opaque). This
change is internal and does not affect users of the Avatica Go client, but maintains best practices and follows the
leading-edge of protobuf usage in Go.

See the list of [bug fixes and new features]({{ site.baseurl }}/docs/go_history.html#v5-4-0)
for more information.