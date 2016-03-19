# Migrating to the "New Environment"

## Status

* `Webhook?`: changes in your repository will trigger jobs in the new jenkins. *See: [Changing Web Hooks](#changing-your-web-hook).*
* `Merged?`: I had to create a branch `geoint` in some repositories; merge to master at your leisure (and remove the `branch: geoint` statement associated with your repository in the Jenkins seed job.
* `Staged?`: is your app running in the new PCF environment?
* `VCAP?`: are you using the VCAP services provided in the new environment?
* `Domain?`: is your app domain agnostic? (i.e. is it still referencing the piazzageo.io domain?)



| App                            | Webhook? | Merged? | Staged? | VCAP? | Domain? |
|--------------------------------|:--------:|:-------:|:-------:|:-----:|:-------:|
| bf-ui                          |          | y       |         |       |         |
| pz-access                      |          |         | y       |       |         |
| pz-discover                    | y        | y       | y       |       |         |
| pz-dispatcher                  |          |         | y       |       |         |
| pz-gateway                     |          |         | y       |       |         |
| pz-ingest                      |          |         | y       |       |         |
| pz-jobcommon                   |          |         | n/a     |       |         |
| pz-jobmanager                  |          |         | y       |       |         |
| pz-logger                      |          | y       | y       |       |         |
| pz-search-lite-metadata-ingest |          |         |         |       |         |
| pz-search-lite-query           |          |         |         |       |         |
| pz-search-metadata-ingest      |          |         |         |       |         |
| pz-search-query                |          |         |         |       |         |
| pz-servicecontroller           |          |         |         |       |         |
| pz-services                    | y        | y       | y       | y     | y       |
| pz-swagger                     | y        | y       | y       |       |         |
| pz-uuidgen                     |          | y       | y       |       |         |
| pz-workflow                    |          | y       | y       |       |         |
| pzclient-sak                   |          | y       | y       |       |         |
| pzsvc-coordinate-conversion    |          | y       | y       |       |         |
| pzsvc-gdaldem                  |          |         | y       |       |         |
| pzsvc-lasinfo                  |          |         | y       |       |         |
| pzsvc-pdal                     |          |         |         |       |         |
| pzsvc-us-phone-number-filter   |          | y       |         |       |         |
| pzsvc-us-geospatial-filter     |          | y       |         |       |         |
| pztest-integration             |          | y       | n/a     |       |         |
| time-lapse-viewer              |          | y       |         |       |         |


## Changing Your Web Hook

1. Navigate to your repository page on github.
1. Click on the **`Settings`** tab.
1. Click on **`Webhooks & Services`** in the menu on the left.
1. Click on the pencil icon (edit) by the **`Jenkins (GitHub plugin)`**
1. Change the **`Jenkins hook url`** to `https://jenkins.devops.geointservices.io/github-webhook/`
1. Click **`Update service`**
![Jenkins Webhook](./img/jenkins-webhook.png)
