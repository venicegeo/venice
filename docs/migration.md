# Migrating to FADE

FADE is the "New Environment" ([More details](devops.md)).

## Status

* **Building?**: all jobs have successfully built in the new jenkins (at least once).
* **VCAP?**: are you using the VCAP services provided in the new environment?
* **Domain?**: is your app domain agnostic? (i.e. is it still referencing the piazzageo.io domain?)


| App                            | Building? | VCAP? | Domain? |
|--------------------------------|:---------:|:-----:|:-------:|
| bf-ui                          | y         |       |         |
| pz-access                      | y         | y     |         |
| pz-discover                    | y         |       |         |
| pz-dispatcher                  | y         | y     |         |
| pz-gateway                     | y         | y     |         |
| pz-ingest                      | y         | y     |         |
| pz-jobcommon                   | y         | n/a   |         |
| pz-jobmanager                  | y         | y     |         |
| pz-logger                      | y         |       | y       |
| pz-search-metadata-ingest      | y         |       |         |
| pz-search-query                | y         |       |         |
| pz-servicecontroller           | y         |       |         |
| pz-services                    | y         | y     | y       |
| pz-swagger                     | y         | n/a   |         |
| pz-uuidgen                     | y         |       | y       |
| pz-workflow                    | y         |       |         |
| pzclient-sak                   | y         |       |         |
| pzsvc-coordinate-conversion    | y         |       |         |
| pzsvc-gdaldem                  | y         |       |         |
| pzsvc-lasinfo                  | y         |       |         |
| pzsvc-pdal                     | y         |       |         |
| pztest-integration             |           | n/a   |         |
| time-lapse-viewer              | y         |       |         |
