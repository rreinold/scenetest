# ClearBlade's SceneTest - Load Testing Tool  

## Overview  

Scenetest provides a simple means of setting up, running, and tearing down tests. The test scripts are written in json (with comments!), and interpreted by the scenetest driver. A single test can spawn multiple “sub tests” in parallel and scenetest provides a rudimentary means of communication and synchronization throughout the test. Tests can be run serially or in parallel.

In general, you can use scenetest to set up a testing environment (systems, users, anything). Next, you can run tests that refer to the previously set up environment. Finally you can teardown a previously set up test environment. We will first look at an example of an actual test. We will ignore setup and teardown for now. Just understand that for your tests, there are prebuilt users, collections, services, libraries, triggers, etc based on the contents of your setup file. Setup and teardown will be discussed in detail later in this tome.

## Concepts

#### Files

1. Setup Configuration (setup.json)

This file contains all assets and their respective configurations.

Ex. A code service with code, and permissions, along with a collection with some data

2. Resource Map (autogenerated)

This is a temporary file used to map all the requested resources with their respective autogenerated ids. It is the output of setup, and input of run.

3. Run Configuration (run.js)

This file contains the set of actions that make up the test scenario itself.

#### Commands

1. `scenetest setup` ingests the Setup Configuration lays down groundwork on which a load test can run. This can be anywhere from an empty system, to a complex system representing an entire IoT solution. This command outputs a Resource Map file, which is ingested by `scenetest run`

2. `scenetest run` ingests the Resource Map file, and uses it to execute a load test. The test is represented by a JSON file containing a series of actions.

3. `scenetest teardown` ingests the Resource Map to destroy all resources used for the test.

## Configure

Clone then run these commands in the **scenetest** folder

```
go get  
go install
```

## Run a Sample Test

1. Run the following command:

```
scenetest setup -info resourceMap.json -platform-url "http://127.0.0.1:9000" -messaging-url "127.0.0.1:1883" setup/novi.json

```
2. You now have a new developer and system containing an entire IoT system. Let's run a series of actions against that system:

```
scenetest run -info resourceMap.json -platform-url "http://127.0.0.1:9000" -messaging-url "127.0.0.1:1883" examples/createItem.json
```
