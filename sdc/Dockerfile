ARG SDC_VERSION=3.21.0
FROM streamsets/datacollector:${SDC_VERSION}

COPY --chown=sdc:sdc ./data /data
RUN "${SDC_DIST}/bin/streamsets" stagelibs -install=streamsets-datacollector-redis-lib,streamsets-datacollector-influxdb_0_9-lib

EXPOSE 18630