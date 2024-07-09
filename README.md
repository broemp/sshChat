# sshChat

This is a Project for my Communication Networks Lecture.

## What does it do?

sshChat creates a ssh Server where one can connect and chat with other people.

## How to build

### Binary

You need Go>=1.22.4

```bash
go mod download # download dependencies
go build .      # build binary
./sshChat       # run binary
```

### Docker

You need docker installed

```bash
docker build . -t [USERNAME]/ssh_chat:[VERSION]

# e.g.
docker build . -t broemp/ssh_chat:dev
docker run -e HOST=127.0.0.1 -e PORT=1337 -p 1337:1337 --name= sshChat -d broemp/ssh_chat:dev
```
