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
