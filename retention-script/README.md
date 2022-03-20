Retention script for Docker Registry
====================================

# How to use

Run directly
```
cd retention-script

go run main.go
```

Build binary (need to set env for different OSes)
```bash
cd retention-script
mkdir bin

env GOOS=darmin GOARCH=amd64 go build -o ./bin/retention ./main.go

env GOOS=linux GOARCH=amd64 go build -o ./bin/retention ./main.go
```

# Ideas
Some ideas for the retention scripts

1. Define a rule list for each service, how many latest tags we will keep
2. Scan through all repo, if the repo in the rule list
3. Get all tags of that repo
4. Delete (total_tags - num_latest_tags)

After deleting the tag, we can wait for the GC to collect the disk space

Otherwise, we can also delete the disk space manually by ourselves, based on the digest of images

# Remaining issues/TODOs

1. Encountering the MANIFEST_UNKNOWN error when deleting image with digest, seems this is a well known issue within the docker-distribution pkg

Reference here

https://github.com/distribution/distribution/issues/1566

https://betterprogramming.pub/cleanup-your-docker-registry-ef0527673e3a

https://github.com/distribution/distribution/issues/1755

2. Improve the script to support json rule, no need to hard code the rules inside the script

3. Build a Docker image then use K8s cronjob to run the retention script

4. We can also try Harbor, which support retention policy by default
