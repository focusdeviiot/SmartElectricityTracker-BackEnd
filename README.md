# docker build
```
docker buildx build --platform linux/amd64,linux/amd64/v2,linux/amd64/v3,linux/arm64,linux/arm/v7,linux/arm/v6 -t focusdeviiot/smart-electricity-tracker-server:1.0.1 -t focusdeviiot/smart-electricity-tracker-server:latest --push .

```