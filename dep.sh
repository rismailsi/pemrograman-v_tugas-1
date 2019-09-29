#!/bin/bash

docker pull registry.softint.com.au/softint/dep:latest

docker run --rm \
    --name=dep-$$ \
    -v ~/.ssh/id_rsa:/root/.ssh/id_rsa \
    -v ~/.ssh/id_rsa.pub:/root/.ssh/id_rsa.pub \
    -v ~/.ssh/known_hosts:/root/.ssh/known_hosts \
    -v $(pwd):/go/src/bitbucket.org/reliefsoftware/service-user \
    -w "/go/src/bitbucket.org/reliefsoftware/service-user" \
    registry.softint.com.au/softint/dep:latest ${*}
