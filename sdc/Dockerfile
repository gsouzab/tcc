ARG SDC_VERSION=3.5.1
FROM streamsets/datacollector:${SDC_VERSION}

RUN "${SDC_DIST}/bin/streamsets" stagelibs -install=streamsets-datacollector-redis-lib,streamsets-datacollector-influxdb_0_9-lib

EXPOSE 18630