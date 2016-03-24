# Migrating to FADE

FADE is the "New Environment" ([More details](devops.md)).

## Status

* **Building?**: all jobs have successfully built in the new jenkins (at least once).
* **Webhook?**: changes in your repository will trigger jobs in the new jenkins. *See: [Changing Web Hooks](#changing-your-web-hook).*
* **Merged?**: I may have created a `geoint` branch in your repo; merge to master and update the Jenkins seed job.
* **Staged?**: is your app running in the new PCF environment?
* **VCAP?**: are you using the VCAP services provided in the new environment?
* **Domain?**: is your app domain agnostic? (i.e. is it still referencing the piazzageo.io domain?)



| App                            | Building? | Webhook? | Merged? | Staged? | VCAP? | Domain? |
|--------------------------------|:---------:|:--------:|:-------:|:-------:|:-----:|:-------:|
| bf-ui                          |           |          | y       |         |       |         |
| pz-access                      | y         | y        | y       | y       | y     |         |
| pz-discover                    | y         | y        | y       | y       |       |         |
| pz-dispatcher                  | y         | y        | y       | y       | y     |         |
| pz-gateway                     | y         | y        | y       | y       | y     |         |
| pz-ingest                      | y         | y        | y       | y       | y     |         |
| pz-jobcommon                   | y         | y        | y       | n/a     | n/a   |         |
| pz-jobmanager                  | y         | y        | y       | y       | y     |         |
| pz-logger                      | y         | y        | y       | y       |       | y       |
| pz-search-metadata-ingest      |           |          |         |         |       |         |
| pz-search-query                |           |          |         |         |       |         |
| pz-servicecontroller           | y         | y        |         | y       |       |         |
| pz-services                    | y         | y        | y       | y       | y     | y       |
| pz-swagger                     | y         | y        | y       | y       | n/a   |         |
| pz-uuidgen                     | y         | y        | y       | y       |       | y       |
| pz-workflow                    | y         |          | y       | y       |       |         |
| pzclient-sak                   | y         |          | y       | y       |       |         |
| pzsvc-coordinate-conversion    |           |          | y       | y       |       |         |
| pzsvc-gdaldem                  |           |          |         | y       |       |         |
| pzsvc-lasinfo                  |           |          |         | y       |       |         |
| pzsvc-pdal                     |           |          |         |         |       |         |
| pzsvc-us-phone-number-filter   |           |          | y       |         |       |         |
| pzsvc-us-geospatial-filter     |           |          | y       |         |       |         |
| pztest-integration             |           |          | y       | n/a     | n/a   |         |
| time-lapse-viewer              |           |          | y       |         |       |         |


## Changing Your Web Hook

1. Navigate to your repository page on github.
1. Click on the **`Settings`** tab.
1. Click on **`Webhooks & Services`** in the menu on the left.
1. Click on the pencil icon (edit) by the **`Jenkins (GitHub plugin)`**
1. Change the **`Jenkins hook url`** to `https://jenkins.devops.geointservices.io/github-webhook/`
1. Click **`Update service`**
![Jenkins Webhook](./img/jenkins-webhook.png)
