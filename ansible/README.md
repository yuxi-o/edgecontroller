# Ansible Deployment

## Overview
[Ansible](https://www.ansible.com) is an infrastructure automation management tool. This project provides Ansible
scripts to get you started in deploying the Controller in your infrastructure. If you're reading this but haven't read
the main documentation yet, read it [here](../README.md). The context of these scripts is closely related to that
documentation.

## Usage

### Configure an ansible inventory

There is an example inventory included called `development` included in this folder.

### Run the playbook

Run the playbook passing in values for the `CONTROLLER_HOST` (the API), and the `UI_HOST`.

In this example, we're specifying the `development` inventory file.  We're using `server1` for both the controller and
the UI.  The additional hosts in our inventory will go unused.  The values passed in for `CONTROLLER_HOST` and
`UI_HOST` must exist in the inventory file.

```sh
ansible-playbook -i development controller-ce.yml \                                                                         2 ↵  03:22 Dur
        -e "controller_host=server1.example.com \
                ui_host=server1.example.com \
                github_token=mysupersecretapikey"

```
