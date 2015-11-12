&tau;
=====
[![Build Status](https://travis-ci.org/jad-b/Torque.svg?branch=master)](https://travis-ci.org/jad-b/Torque)

__Torque__ is a platform for collecting, analyzing, and acting upon personal data.
It is geared towards personal betterment, not treatment of serious medical conditions. Although, as they say, life is a terminal condition to have, so you never know what will come out of a side project.

This data (will) include everything from workouts to basic biomarkers to blood chemistry to psychological profiles.


Table of Contents
* [Deployment](deploy/README.md)
* [Code Structure](#structure)
* [Design](docs/Design.md)

## Structure
All general-purpose libraries and interface definitions are defined in the
top-level `torque` directory. Resource definitions of a related kind, such as
biometrics or workout data, belong in their own sub-directory package.

