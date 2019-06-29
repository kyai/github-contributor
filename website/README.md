#### nginx

```
server {
    listen      80;
    server_name githubcontributor.com;

    root        /www/github-contributor/website;

    location / {
        index   index.html;
    }

    location ~* /data/(.+)/(.+) {
        resolver    8.8.8.8;
        proxy_pass  https://github.com/$1/$2/graphs/contributors-data;
    }
}
```