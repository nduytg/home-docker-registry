Write a script that implements Docker Registry retention. In the registry, we might have many repositories with different tags
* serviceA:tag1, serviceA:tag2, serviceA:v7, service:latest
* serviceB:v3, serviceB:v1, etc.

The retention job/script should be able to work with different retention settings: For example:
* Retain only the last 5 images of serviceA
* Retain only the last 3 images of serviceB
* Retain only the last X images of serviceY
 
Feel free to choose any language and method to periodically run the retention job.