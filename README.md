# metadata-reaper-for-kafka
kafka metadata reaper


## Overview
Kafka metadata reaper is a service written in golang who will help you to delete the metadata created in your kafka cluster. This service is going to create a admin client against your kafka cluster and list all topics.

## Getting Started

### Build & Run
The repo provide a makefile with some useful instruccion. You can start a kafka cluster with zookeeper using the docker-compose file in this project. This will deploy a simple local cluster in which you can test your service.
```bash
    make start-environ
```

Nevertheless, you can build the service and run the binary. Next command will generate a binary in `bin/kafka-metadata-reaper` who you can run locally.
```bash
    make build
```

### Run kafka metadata reaper instance
Once you have the kafka cluster up and running you can start your kafka-metadata-reaper instance:
- Topics: List empty topics ready for deletion (dry-run mode).
- Topics (no dry-run): List and delete empty topics.
You can run each one using the next commands:
```bash
    make run-topics-reaper
    make run-topics-reaper-no-dryrun
```

## Contributing

The metadata-reaper-for-kafka project team welcomes contributions from the community. Before you start working with metadata-reaper-for-kafka, please
read our [Developer Certificate of Origin](https://cla.vmware.com/dco). All contributions to this repository must be
signed as described on that page. Your signature certifies that you wrote the patch or have the right to pass it on
as an open-source patch. For more detailed information, refer to [CONTRIBUTING_DCO.md](CONTRIBUTING_DCO.md).

## License

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.  IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.