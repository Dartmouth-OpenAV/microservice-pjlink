# openav microservice pjlink

OpenAV microservice for PJLink.  Uses the microservice framework and PJLink to communicate with projectors.

This microservice operates and controls projectors that adhere to the PJLink communication standard.

## Endpoint Notes

Some devices throw errors when trying to use `audiomute` or `videomute` independently. For these devices, use the `audioandvideomute` endpoint instead.

##

[JBMIA PJLink](https://pjlink.jbmia.or.jp/english/)

[SONY VPL-FHZ85](https://pro.sony/ue_US/products/laser-projectors/vpl-fhz85)

[Microservice testing documentation](https://github.com/Dartmouth-OpenAV/documentation/blob/main/curl_test_readme.md)

![](https://github.com/Dartmouth-OpenAV/microservice-pjlink/blob/main/diagram.png)

![](https://github.com/Dartmouth-OpenAV/microservice-pjlink/blob/main/front.png)

![](https://github.com/Dartmouth-OpenAV/microservice-pjlink/blob/main/rear.png)
