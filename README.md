# Shelly H&T Prometheus Exporter

This is a Prometheus exporter for [Shelly H&T sensor IoT devices](https://shellyparts.de/en/products/shelly-h-t).
The exporter supports multiple Shelly H&T sensors which can send their data to the exporter via the `Actions` integration.


## Building and Running the Exporter

To build the project, run:

```
go build
```

Run the built executable using:

```
./shelly-ht-exporter
```

This will start the web server.
By default, it listens on `127.0.0.1:8090`.

To listen on a different address, use the `SHELLY_HT_EXPORTER_ADDR` environment variable:
```
export SHELLY_HT_EXPORTER_ADDR=":8000"
./shelly-ht-exporter
```


## Configure Shelly H&T to Send Data

To have your Shelly H&T sensor send data to the Prometheus exporter, you should configure it to send HTTP requests to the exporter.

1. **Access Shelly Device Configuration**:
   - Ensure your Shelly H&T sensor is connected to the network and accessible through its web interface.

2. **Set Up Actions**:
   - Go to the settings of Shelly H&T and go to the tab for configuring actions
   - Add a new action, specifying the URL of your Prometheus exporter. This would be something like `http://<your_exporter_ip>:8090/`.

3. **Testing**:
   - Once set up, the Shelly H&T should start sending requests to the exporter every five minutes
   - You can check the `/metrics` endpoint of the exporter to see if the data has been updated.

## Map Sensors to Human Readaable IDs

The Shelly H&T sensors all have a unique identifier like `shellyht-CC5AS8`.
Unfortunately, this is not very nice to remember if you have several sensors.
To make your life easier, you can use the `SHELLY_HT_EXPORTER_NAME_MAP` environment variable to map those IDs to nice, human readable names.

```
export SHELLY_HT_EXPORTER_NAME_MAP='{"shellyht-CC5AS8":"living room"}'
./shelly-ht-exporter
```

The value of that variable should be a JSON object specifying a map of IDs and names.
The names should be unique.
They will replace the actual ID in the metrics.


## Metrics

The exporter provides a `/metrics` endpoint.
You can use cURL to test it and get an overview of the available data:

```
❯ curl http://localhost:8090/metrics

# HELP count_updated The number of updates to a sensor value
# TYPE count_updated counter
count_updated{sensor="shellyht-221"} 1
count_updated{sensor="shellyht-CA62E1"} 2
...
# HELP humidity The measured humidity
# TYPE humidity gauge
humidity{sensor="shellyht-221"} 52.4
humidity{sensor="shellyht-CA62E1"} 54.2
# HELP last_updated_time_seconds Timestamp of the last update in seconds since epoch.
# TYPE last_updated_time_seconds gauge
last_updated_time_seconds{sensor="shellyht-221"} 1.731776654e+09
last_updated_time_seconds{sensor="shellyht-CA62E1"} 1.73177662e+09
...
# HELP temperature The measured temperature
# TYPE temperature gauge
temperature{sensor="shellyht-221"} 20.7
temperature{sensor="shellyht-CA62E1"} 21.5

```
