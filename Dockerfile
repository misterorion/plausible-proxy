ARG caddy_image
FROM $caddy_image
COPY Caddyfile /etc/caddy/Caddyfile