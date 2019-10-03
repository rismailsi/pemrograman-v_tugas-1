# pemrograman-v_tugas-1

## Run Environment

Make sure you have installed docker and docker-compose

Install the Go vendors.  Should create a `vendor` directory in the projects root.
***Note:*** this mounts ~/.ssh into /root/.ssh in the dep container,
this allows dep to pull private repos using your ssh keys.

```shell
./dep.sh ensure -v
```

Pull and build the development environment images
```shell
docker-compose pull && docker-compose build --pull
```

Bring up the environment
```shell
docker-compose up
```

Browse to http://localhost to test the app

Browse to http://localhost:8081 to access pgweb

