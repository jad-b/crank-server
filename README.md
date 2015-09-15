&tau;
=====
[![Build Status](https://travis-ci.org/jad-b/Torque.svg?branch=master)](https://travis-ci.org/jad-b/Torque)

__Torque__ is a platform for collecting, analyzing, and acting upon personal data.
It is geared towards personal betterment, not treatment of serious medical conditions. Although, as they say, life is a terminal condition to have, so you never know what will come out of a side project.

This data (will) include everything from workouts to basic biomarkers to blood chemistry to psychological profiles.


Table of Contents
* [Deployment](deploy/README.md)
* [Code Structure](#structure)

## Structure
All general-purpose libraries and interface definitions are defined in the
top-level `torque` directory. Resource definitions of a related kind, such as
biometrics or workout data, belong in their own sub-directory package.

### Why do your `structs` have so many methods?
After much debate, it was decided that attempting to force a class-based
inheritance hierarchy onto Go was a step in the wrong direction. What this
means is that there are Actors in Torque (or Resources, or Objects - whatever
term resonates with you) that implement potentially _many_ interfaces. For
instance, a basic metric like Bodyweight acts as a Database Resource, a REST
API Resource, and a CLI client. All methods are defined by this resource.

The alternative was to provide a basic struct definition - a Bodyweight record
has a weight, timestamp, comment, and user associated with it - and
embed/compose that within other structs. But this didn't really show any
advantages, down the road, for when we'd need to refactor (and you *always*
need to refactor). What would be good about having to look in five different
places for how a certain resource behaved in different situations? This would
be a potential solution to working _around_ a constraint, such as not having
access to modify the underlying object directly; say it's being imported from
a public library. But when you can delegate, you delegate - let the Actor tell
you how _it_ behaves. Anyway, hopefully this makes sense in the morning.

