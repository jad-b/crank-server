redteam
=======

Red Team represents the real enemy - the user. It is inherently system-level, as it wants to interact with the entire
system from the point of view of a specific user.

Red Team scripts should *always*:
* Start with a client library
* Rely on the entire stack being in place
* *Never* use code only available to superusers and the internal Torque stack
* *Never* rely on being a superuser - if it needs superuser privilegs to setup, do so in a setup script.

Once the `redteam` user has been created, it should be interacted with as any other user account.

Eventually, it should be capable of running random, scheduled events, much like
Netflix's Chaos Monkey series.
