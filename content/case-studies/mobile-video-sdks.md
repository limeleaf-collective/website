+++
title = "Networking Case Study: End-to-end HTTP/2 Traffic"
description = ""
weight = 2

[extra]
feature_photo = "images/services/prod-dev.svg"
feature_photo_alt = "Product development illustration"
+++

**Problem:** Modern browsers support HTTP/2 for fast delivery, but edge nodes downgraded traffic to HTTP/1.1.

**Solution:** Configured network edge nodes to retain HTTP/2 traffic.

**Results:** Browser to origin server has full HTTP/2 support; gRPC support for free; bi-directional streaming.