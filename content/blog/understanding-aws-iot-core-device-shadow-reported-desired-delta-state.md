+++
title = "Understanding AWS IoT Core Device Shadow: Reported, Desired, and Delta State"
date = 2025-05-29
draft = false
authors = ["Blain Smith"]

[taxonomies]
tags = ["Engineering", "IoT"]

[extra]
feature_photo = ""
feature_photo_alt = ""
+++

AWS IoT Core Device Shadow is a powerful feature that allows connected devices and cloud applications to synchronize state seamlessly. This feature maintains a *virtual representation* of each device's state in the cloud, even when the device is offline. This blog post will explain how the device shadow works, focusing on the `reported`, `desired`, and `delta` states—and when devices should report their state or respond to changes.

<!-- more -->

## Device Shadow Document Structure

The shadow document is a JSON structure that contains three primary sections:

* **`desired`**: What the cloud application wants the device’s state to be.
* **`reported`**: What the device says its current state is.
* **`delta`**: What’s different between `desired` and `reported`.

Here's a typical shadow document:

```json
{
  "state": {
    "desired": {
      "light_color": "0xfff",
      "temp": 13
    },
    "reported": {
      "light_color": "0x0f0",
      "temp": 21
    },
    "delta": {
      "light_color": "0xfff",
      "temp": 13
    }
  },
  "metadata": {
    ...
  },
  "version": 13,
  "timestamp": 1629999999
}
```

## What Happens When the Device Connects

When a device connects to AWS IoT Core, **it should publish its current state** to the `reported` section of the device shadow. This is essential because:

* It tells the cloud what the device is currently doing.
* It resets any outdated `delta` differences.
* It prevents redundant commands if the device is already in the desired state.

### Reporting State on First Connect

```json
{
  "state": {
    "reported": {
      "light_color": "0x0f0",
      "temp": 21
    }
  }
}
```

The device publishes this JSON to the shadow update topic `$aws/things/device123/shadow/update`. After this, AWS IoT Core compares the new `reported` state with the existing `desired` state and computes the `delta`. A `delta` message is generated **only when** there is a difference between `desired` and `reported` states.

### Setting a New Desired State

A cloud application wants the light color to be white and the temperature to be 13. It publishes this JSON:

```json
{
  "state": {
    "desired": {
      "light_color": "0xfff",
      "temp": 13
    }
  }
}
```

If the device has previously reported:

```json
{
  "state": {
    "reported": {
      "light_color": "0x0f0",
      "temp": 21
    }
  }
}
```

Then the `delta` will be:

```json
{
  "state": {
    "light_color": "0xfff",
    "temp": 13
  }
}
```

This `delta` is published to a specific MQTT topic for the device (e.g., `$aws/things/device123/shadow/update/delta`). The device listens for this topic to know what has changed and should take action accordingly.

## Synchronizing State

Once the device processes the `delta`—for example, it changes the temperature to 13 and it changes the light color to white—it must update its `reported` state again:

```json
{
  "state": {
    "reported": {
      "light_color": "0xfff",
      "temp": 13
    }
  }
}
```

When `reported` matches `desired`, the `delta` disappears:

```json
{
  "state": {
    "desired": {
      "light_color": "0xfff",
      "temp": 13
    },
    "reported": {
      "light_color": "0xfff",
      "temp": 13
    }
  }
}
```

No further `delta` messages are sent unless the desired state changes again.

## Sequence Diagram of State Flow

```
+-------------------+          +---------------------+         +---------------------+
|    Cloud App      |          |   AWS IoT Core      |         |     IoT Device      |
+-------------------+          +---------------------+         +---------------------+
         |                               |                              | <-- Subscribed to delta
         |                               |                              |
         | -- Set desired state -------> |                              |
         |   { "state": {                |                              |
         |       "desired": { ... } } }  |                              |
         |                               |                              |
         |                               | -- Delta calculated ------>  |
         |                               |   { "state": { delta... } }  |
         |                               |                              | -- Takes action to apply changes-+
         |                               |                              | <--------------------------------+
         |                               |                              |
         |                               |                              |
         |                               |                              |
         |                               | <-- Reports new state ------ |
         |                               |  { "state": { reported... } }|
         |                               |                              |
         |                               | <-- Delta cleared if matched |
         |                               |                              |
         |                               |                              |
```

### Legend:

* **Cloud App**: Publishes `desired` state changes to the device shadow on `$aws/things/device123/shadow/update`.
* **AWS IoT Core**:
  * Stores `desired` and `reported`.
  * Computes `delta` when they differ.
  * Sends `delta` messages to the device on `$aws/things/device123/shadow/update/delta`.
* **IoT Device**:
  * Subscribes to `delta` topic `$aws/things/device123/shadow/update/delta`.
  * Acts on `delta` message.
  * Publishes new `reported` state after making changes on `$aws/things/device123/shadow/update`.

### Example Data Flow

1. Cloud sets:

   ```json
   {
     "state": {
       "desired": {
        "light_color": "0xfff",
        "temp": 13
       }
     }
   }
   ```

2. Delta message sent to device:

   ```json
   {
     "state": {
        "light_color": "0xfff",
        "temp": 13
     }
   }
   ```

3. Device updates:

   ```json
   {
     "state": {
       "reported": {
          "light_color": "0xfff",
          "temp": 13
       }
     }
   }
   ```

Once `reported` **equals** `desired`, the `delta` is removed, and the shadow reflects synchronization.

## Best Practices

1. **Report State on Connect**: Devices should always publish their current state upon (re)connection to AWS IoT Core.
2. **Listen to Delta Messages**: Devices should subscribe to the delta topic and act on any changes.
3. **Acknowledge Changes**: After acting on a delta, devices should update their `reported` state to confirm the change.
4. **Avoid Spamming Desired**: Cloud applications should update `desired` only when a change is needed—not continuously—to prevent unnecessary delta messages.
