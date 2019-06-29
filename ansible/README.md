# Ansible Deployment

## Overview
[Ansible](https://www.ansible.com) is an infrastructure automation management tool. This project provides Ansible
scripts to get you started in deploying the Controller in your infrastructure. If you're reading this but haven't read
the main documentation yet, read it [here](../README.md). The context of these scripts is closely related to that
documentation.

## Host Requirements

- CentOS (7.6.1810 or higher)
- Connection to the Internet
- Minimum 2 GB RAM (for building the Docker images)

## Usage

### Configure an ansible inventory

There is an example inventory included called `development` included in this folder.

See the documentation from Ansible on [working with inventory files][inv].
Depending on your environment, you may need to login as a super-user and
[escalate your privileges][escalate].  For example:

    server1.example.com ansible_user=centos ansible_become=yes

[inv]: https://docs.ansible.com/ansible/latest/user_guide/intro_inventory.html
[escalate]: https://docs.ansible.com/ansible/latest/user_guide/become.html

### Run the playbook

Run the playbook passing in values for the `CONTROLLER_HOST` (the API), and the `UI_HOST`.

In this example, we're specifying the `development` inventory file.  We're using `server1` for both the controller and
the UI.  The additional hosts in our inventory will go unused.  The values passed in for `CONTROLLER_HOST` and
`UI_HOST` must exist in the inventory file.

```sh
ansible-playbook -i development controller-ce.yml \
        -e "controller_host=server1.example.com \
                ui_host=server1.example.com \
                github_token=mysupersecretapikey \
                version=0.0.57"
```

In the example, `version=0.0.57` indicates that we wish to install version
`0.0.57` of the controller.  If unspecified, the default will be to use the
`master` branch of the controller.

You will find the controller installed to `/opt/controller` on the host by
default.  You can configure this with the `controller_path` variable.

Credentials for the admin user of the controller as well as the root MySQL user
will be placed in a folder named `credentials/` from where you've run
`ansible-playbook`.
