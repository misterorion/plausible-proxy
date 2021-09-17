# plausible-proxy

Service to reverse proxy analytics to plausible.io

Builds and runs on Google Cloud.

Update paths on load balancer to point at the backend Cloud Run service:

```
/js/script.js
/api/event
```