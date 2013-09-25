OpsMeter
========

A small webserver and arduino program that receives information to display and then displays it on an RGB bar. 

Install that udev rule to /etc/udev/rules.d/72-micro-devel.rules to ensure that the device always shows up as the 
same comm port. Otherwise RESET commands will cause issues with dropped comm ports, as the server will have a harder time finding the comm port. Currently won't be able to anyways however.
