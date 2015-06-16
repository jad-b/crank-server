# Deployment
Uses Ansible to provision and orchestrate deployments. Docker is used to
containerize running services, as well as performing one-off tasks, such
as making a backup of a database.

Deployment codes lives under [deploy/](deploy/). All build commands are
executed by the `robot` user, which is set up by the `robot` role.

## Getting Started
1. Configure your inventory.
1. Add the robot user's private key to your ssh-agent with:
    `ssh-add ~/.ssh/robot_id_rsa`
1. Provision a base system with:
    `ansible-playbook base.yml`

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

__To Implement__
* Binding Docker daemon's to a private/public IP
* Setting up TLS for client/host Docker interactions
    * Can docker-machine handle this?
* Setting up Swarm

## Security
Ideally, the following policies should be in place:
* All requests go through Nginx, which enforces HTTPS
* Nginx performs SSL termination, and load-balances the requests out to the
  backend serveres
* Servers _only_ accept incoming HTTP connections on port 80 from their Nginx
  load-balancer.
* Highest level of TLS is enforced. Go says it "partially supports" TLS1.2
