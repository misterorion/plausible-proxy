# plausible-proxy

Service to reverse proxy analytics to plausible.io

Builds with Cloud Build and deploys to Cloud Run.

Add the paths below to your load balancer and set them to point at the Cloud Run NEG:

```
/js/script.js
/api/event
```

Then add code to your site's `<head>`

```
<script defer data-domain="mysite.com" src="/js/script.js"></script>
```
