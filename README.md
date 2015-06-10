Crank
=====

A RESTful Crank API server (and more), written in Go.

Table of Contents
* [Deployment](#deployment)

## Deployment
Uses Ansible to provision and orchestrate deployments. Docker is used to
containerize running services, as well as performing one-off tasks, such
as making a backup of a database.

Deployment codes lives under [deploy/](deploy/). All build commands are
executed by the `robot` user, which is set up by the `robot` role.
