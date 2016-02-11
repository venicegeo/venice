# VeniceGEO Devops

## Overview

### Infrastructure/Configuration in Code

#### Ephemeral
- Amazon physical resources (ec2 instances).
- These resources will come and go.
- *Should* only be created as a result of an automated process.

#### Runtime
- Application source code.
- AWS Cloudformation templates.
  - All resources defined in a single JSON file, within SCM.
  - Venice cloudformation templates: [canal](https://github.com/venicegeo/canal) _(private repo)_

#### Persistent
- Machine configuration.
- Baked into Amazon Machine Images (AMI) using Chef.
- Venice AMI build process: [gondola](https://github.com/venicegeo/gondola) _(private repo)_


### CI/CD

#### [Jenkins](http://jenkins.piazzageo.io)
- One script defines all build pipelines: [jenkins](https://github.com/venicegeo/jenkins)
- Testing and building artifacts; automated delivery to CloudFoundry.
  - Feedback published in [slack](https://venicegeo.slack.com).
  - Artifacts are currently being stored in s3, but we're move to Nexus soon.
  - build and testing tools baked into the jenkins machine image (via [gondola](https://github.com/venicegeo/gondola)).

![Jenkins Build Dashboard](./img/jenkins-dashboard.png)

![Jenkins Build Pipeline](./img/jenkins-pipeline.png)

#### [CloudFoundry](http://login.cf.piazzageo.io)
- Running Piazza services:
  - [pz-discover](https://github.com/venicegeo/pz-discover)
  - [pz-logger](https://github.com/venicegeo/pz-logger)
  - [pz-alerter](https://github.com/venicegeo/pz-alerter)
  - [pz-uuidgen](https://github.com/venicegeo/pz-uuidgen)
  - [pz-jobmanager](https://github.com/venicegeo/pz-jobmanager)
  - [pz-dispatcher](https://github.com/venicegeo/pz-jobmanager)
  - [pz-gateway](https://github.com/venicegeo/pz-gateway)
  - [pz-servicecontroller](https://github.com/venicegeo/pz-servicecontroller)
  - [pzsvc-gdaldem](https://github.com/venicegeo/pz-gdaldem)
  - [pzsvc-lasinfo](https://github.com/venicegeo/pz-lasinfo)

#### Shared Services
- _Note: these services are currently running in the larger AWS environment, but they are being migrated to CloudFoundry._
  - elasticsearch 
  - geoserver
  - kafka
  - mongodb
  - postgresql (with PostGIS)
  - zookeeper

- Not migrating to CloudFoundry, but available nonetheless:
  - BOSH
  - CloudFoundry
  - jenkins
  - swagger

## Status Update: 11 Feb 2016

_The State of VeniceGEO Devops_

### Past

#### First pass infrastructure complete.

##### Piazza specific:

  * pz-* core components _cf_
  * pzsvc-* services _cf_
  * Swagger _aws, cf2_
  * Kafka _aws_
  * Zookeeper _aws_
  * MongoDB _aws_
  * PostGIS _aws_
  * Elasticsearch _aws_
  * Geoserver _aws_

##### Temporary Infrastructure components:
  * Jenkins _aws_
  * OS Cloud Foundry (`*.cf.piazzageo.io`) _aws_

### Present

#### Migration to GEOINT

  * Pivotal Cloud Foundry (`*.cf2.piazzageo.io`) || GEOINT's PCF
  * Deploy Piazza's deps to PCF (e.g. kafka, zk, etc) and make case for non-PCF services.
  * Build piazza components on GEOINT Jenkins.
    * Artifacts in NEXUS
    * Blue/Green deploys to GEOINT PCF.

#### Continue to stand up new services and support.

  * GeoSHAPE
  * logstash?

### Future

#### Security

  * HP Fortify
  * S3 Bucket policies
  * Secrets sharing
  * SSL
