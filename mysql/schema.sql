-- Copyright 2019 Smart-Edge.com, Inc. All rights reserved.
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--     http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

DROP DATABASE IF EXISTS controller_ce;

CREATE DATABASE controller_ce;

USE controller_ce

-- -------------
-- Entity tables
-- -------------

CREATE TABLE nodes (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    entity JSON
);

CREATE TABLE container_apps (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    entity JSON
);

CREATE TABLE vm_apps (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    entity JSON
);

CREATE TABLE container_vnfs (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    entity JSON
);

CREATE TABLE vm_vnfs (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    entity JSON
);

CREATE TABLE traffic_policies (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    entity JSON
);

CREATE TABLE dns_configs (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    entity JSON
);

CREATE TABLE dns_container_app_aliases (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    container_app_id VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.container_app_id') STORED,
    entity JSON,
    FOREIGN KEY (container_app_id) REFERENCES container_apps(id)
);

CREATE TABLE dns_vm_app_aliases (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    vm_app_id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.vm_app_id') STORED,
    entity JSON,
    FOREIGN KEY (vm_app_id) REFERENCES vm_apps(id)
);

CREATE TABLE dns_container_vnf_aliases (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    container_vnf_id VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.container_vnf_id') STORED,
    entity JSON,
    FOREIGN KEY (container_vnf_id) REFERENCES container_vnfs(id)
);

CREATE TABLE dns_vm_vnf_aliases (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    vm_vnf_id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.vm_vnf_id') STORED,
    entity JSON,
    FOREIGN KEY (vm_vnf_id) REFERENCES vm_vnfs(id)
);

-- -------------------
-- Primary join tables
-- -------------------

-- These tables join two entity tables.

-- dns_configs x dns_container_app_aliases
CREATE TABLE dns_configs_dns_container_app_aliases (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    dns_config_id  VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.dns_config_id') STORED,
    dns_container_app_alias_id  VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.dns_container_app_alias_id') STORED,
    entity JSON,
    FOREIGN KEY (dns_config_id) REFERENCES dns_configs(id),
    FOREIGN KEY (dns_container_app_alias_id) REFERENCES
        dns_container_app_aliases(id),
    UNIQUE KEY (dns_config_id, dns_container_app_alias_id)
);

-- dns_configs x dns_vm_app_aliases
CREATE TABLE dns_configs_dns_vm_app_aliases (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    dns_config_id  VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.dns_config_id') STORED,
    dns_vm_app_alias_id  VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.dns_vm_app_alias_id') STORED,
    entity JSON,
    FOREIGN KEY (dns_config_id) REFERENCES dns_configs(id),
    FOREIGN KEY (dns_vm_app_alias_id) REFERENCES dns_vm_app_aliases(id),
    UNIQUE KEY (dns_config_id, dns_vm_app_alias_id)
);

-- dns_configs x dns_container_vnf_aliases
CREATE TABLE dns_configs_dns_container_vnf_aliases (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    dns_config_id  VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.dns_config_id') STORED,
    dns_container_vnf_alias_id  VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.dns_container_vnf_alias_id') STORED,
    entity JSON,
    FOREIGN KEY (dns_config_id) REFERENCES dns_configs(id),
    FOREIGN KEY (dns_container_vnf_alias_id) REFERENCES
        dns_container_vnf_aliases(id),
    UNIQUE KEY (dns_config_id, dns_container_vnf_alias_id)
);

-- dns_configs x dns_vm_vnf_aliases
CREATE TABLE dns_configs_dns_vm_vnf_aliases (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    dns_config_id  VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.dns_config_id') STORED,
    dns_vm_vnf_alias_id  VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.dns_vm_vnf_alias_id') STORED,
    entity JSON,
    FOREIGN KEY (dns_config_id) REFERENCES dns_configs(id),
    FOREIGN KEY (dns_vm_vnf_alias_id) REFERENCES
        dns_vm_vnf_aliases(id),
    UNIQUE KEY (dns_config_id, dns_vm_vnf_alias_id)
);

-- nodes x container_apps
CREATE TABLE nodes_container_apps (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    node_id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.node_id') STORED,
    container_app_id VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.container_app_id') STORED,
    entity JSON,
    FOREIGN KEY (node_id) REFERENCES nodes(id),
    FOREIGN KEY (container_app_id) REFERENCES container_apps(id),
    UNIQUE KEY (node_id, container_app_id)
);

-- nodes x vm_apps
CREATE TABLE nodes_vm_apps (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    node_id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.node_id') STORED,
    vm_app_id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.vm_app_id') STORED,
    entity JSON,
    FOREIGN KEY (node_id) REFERENCES nodes(id),
    FOREIGN KEY (vm_app_id) REFERENCES vm_apps(id),
    UNIQUE KEY (node_id, vm_app_id)
);

-- nodes x container_vnfs
CREATE TABLE nodes_container_vnfs (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    node_id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.node_id') STORED,
    container_vnf_id VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.container_vnf_id') STORED,
    entity JSON,
    FOREIGN KEY (node_id) REFERENCES nodes(id),
    FOREIGN KEY (container_vnf_id) REFERENCES container_vnfs(id),
    UNIQUE KEY (node_id, container_vnf_id)
);

-- nodes x vm_vnfs
CREATE TABLE nodes_vm_vnfs (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    node_id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.node_id') STORED,
    vm_vnf_id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.vm_vnf_id') STORED,
    entity JSON,
    FOREIGN KEY (node_id) REFERENCES nodes(id),
    FOREIGN KEY (vm_vnf_id) REFERENCES vm_vnfs(id),
    UNIQUE KEY (node_id, vm_vnf_id)
);

-- nodes x dns_configs
CREATE TABLE nodes_dns_configs (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    node_id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.node_id') STORED,
    dns_config_id VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.dns_config_id') STORED,
    entity JSON,
    FOREIGN KEY (node_id) REFERENCES nodes(id),
    FOREIGN KEY (dns_config_id) REFERENCES dns_configs(id),
    UNIQUE KEY (node_id)
);

-- ---------------------
-- Secondary join tables
-- ---------------------

-- These tables join an entity table to a primary join table.

-- nodes_container_apps x traffic_policies
CREATE TABLE nodes_container_apps_traffic_policies (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    nodes_container_apps_id VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.nodes_container_apps_id') STORED,
    traffic_policy_id VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.traffic_policy_id') STORED,
    entity JSON,
    FOREIGN KEY (nodes_container_apps_id) REFERENCES nodes_container_apps(id),
    FOREIGN KEY (traffic_policy_id) REFERENCES traffic_policies(id),
    UNIQUE KEY (nodes_container_apps_id, traffic_policy_id)
);

-- nodes_vm_apps x traffic_policies
CREATE TABLE nodes_vm_apps_traffic_policies (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    nodes_vm_apps_id VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.nodes_vm_apps_id') STORED,
    traffic_policy_id VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.traffic_policy_id') STORED,
    entity JSON,
    FOREIGN KEY (nodes_vm_apps_id) REFERENCES nodes_vm_apps(id),
    FOREIGN KEY (traffic_policy_id) REFERENCES traffic_policies(id),
    UNIQUE KEY (nodes_vm_apps_id, traffic_policy_id)
);
