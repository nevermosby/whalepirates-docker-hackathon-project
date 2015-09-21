#!/bin/sh

docker run -d -p 80:4000 \
	-v /home/ubuntu/DockerHackDay-201509/web:/usr/src/app \
	--name hackathon-web david/hackathon-web:0.1
