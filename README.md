# Raspberry Pi Thermostat

This program was forked from https://github.com/sjenning/rpi-thermostat

I made changes to handle 240v base board using solid state relay board and 240v contactor. It supports 3 channels but can be upgraded to many more if needed. Each channel generates a pwm based on a simple proportionnal controller.

I also changed the temperature sensor to use a 1-wire DS18B20+.

--------------

This project consists of a simple Go program and web UI for using a Raspberry Pi as an web controllable thermostat.

Please refer to this blog post for more information:

https://www.variantweb.net/blog/building-a-thermostat-with-the-raspberry-pi/
