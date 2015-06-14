# Deployment
Uses Ansible to provision and orchestrate deployments. Docker is used to
containerize running services, as well as performing one-off tasks, such
as making a backup of a database.

Deployment codes lives under [deploy/](deploy/). All build commands are
executed by the `robot` user, which is set up by the `robot` role.

## Secrets
You will be expected to have your vault password in a `.vault_pw.txt` file, as
configured by `ansible.cfg`.
You will need to run the `robot.yml` playbook at least once as a `sudo` user on the
remote machine. From there on out, you should be able to use the `robot` user,
which is set in `ansible.cfg` as well. For the first bootstrapping, you can use
the `-u $USER` option.

## Docker
We're using [Docker's bootstrap script](https://get.docker.io/) for setting
up Docker on our hosts. That can be done with this command:
    `ansible-playbook docker.yml`

## Security
Ideally, the following policies should be in place:
* All requests go through Nginx, which enforces HTTPS
* Nginx performs SSL termination, and load-balances the requests out to the
  backend serveres
* Servers _only_ accept incoming HTTP connections on port 80 from their Nginx
  load-balancer.
