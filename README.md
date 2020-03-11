# Porfoliotracker message consumer golang microservice

A simple golang message consumer microservice. 


## Usage 

```bash
# cd to project directory and build executable
$ go build -o microservice .

```

## Docker build

```bash
docker build -t <your-registry-id>/myportfolio-message-consumer:1.12.5 .

```

## Executing tests
As the consumer relies on an event to trigger (from the kafka queue), it was easier to setup a standalone kafka docker image
and do some E2E testing. The connectors are basically setup in the consumer_test.go file (connectors.go and main.go are ignored)

```bash
# The first step is to execute the docker-compose.yml file. This will start a Kafka and Zookeeper setup
# execute the command 'docker ps -a' to ensure the images are running
docker-compose up

# once Kafka and Zookeeper are verified running, execute the following commands
docker ps -a 
# find the Container ID for Kafka
# execute these commands 
docker exec -it <containerid-kafka> bash
# this will write a message to the queue
root@ccfe16688f2d:/$ echo "test-XYZ" | kafka-console-producer --broker-list kafka:29092 --topic analytics
# in a separate terminal execute the following
# clear the cache - this is optional
go clean -testcache
go test -v schema.go consumer.go consumer_test.go -coverprofile tests/results/cover.out
# To view the code coverage in html execute the following command
go tool cover -html=tests/results/cover.out -o tests/results/cover.html
# The output from the test shouild pass and have about 80% coverage
coverage: 80.0% of statements
ok  	command-line-arguments	0.019s	coverage: 80.0% of statements
# The next step is optional (for code analysis - style and linting, coverage, bugs, vulnerabilities, code smell etc)
# run sonarqube scanner (assuming sonarqube server is running)
# NB the SonarQube host and login will differ - please update it accordingly 
 ~/Programs/sonar-scanner-3.3.0.1492-linux/bin/sonar-scanner  -Dsonar.projectKey=myportfolio-message-consumer  -Dsonar.sources=.   -Dsonar.host.url=http://localhost:9009   -Dsonar.login=3b172e408d048820bc6a633b1c3f0097523e89f4 -Dsonar.go.coverage.reportPaths=tests/results/cover.out -Dsonar.exclusions=vendor/**,*_test.go,main.go,connectors.go,tests/**

```
