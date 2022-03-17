# The task

## Main task

Write the Kubernetes deployment manifest to run Docker Registry in Kubernetes with at least the following resources:
[] deployment
[] service
[] persistent volume claim
[] garbage collect cron job
[] ingress
[] secret (if needed). 
[] configmap (added)
[] self-signed ssl (added)

Make sure Docker Registry uses Redis for cache and look into the possibility to run 2 replicas of Docker Registry for redundancy.

Record any issues you encounter. You’re also welcome to list improvement ideas for the service that are outside the scope of this assignment.

## Optional tasks
Write a script that implements Docker Registry retention. In the registry, we might have many repositories with different tags
* serviceA:tag1, serviceA:tag2, serviceA:v7, service:latest
* serviceB:v3, serviceB:v1, etc.

The retention job/script should be able to work with different retention settings: For example:
* Retain only the last 5 images of serviceA
* Retain only the last 3 images of serviceB
* Retain only the last X images of serviceY
 
Feel free to choose any language and method to periodically run the retention job.

## Solution
Please return your finished solution as a .zip file or GitHub repo to dechan@33talent.com 


