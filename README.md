Contains performance benchmark tool for Mongo
Run ```make stress``` to start the stress tool.
Change the Mongo IP/Port in options.go and then ```make build```

Start the stress tool from multiple hosts and start stressing Mongo parallelly.

Check the ```mongo-sharding-config``` file for sharding configuration for Mongo.
```stat.sh``` should be run in all the shard replicas (and config replicas) to get the stat details.


