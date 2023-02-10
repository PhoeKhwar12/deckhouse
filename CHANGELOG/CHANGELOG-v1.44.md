# Changelog v1.44

## [MALFORMED]


 - #3505 unknown section "global"
 - #3568 unknown section "global"
 - #3633 invalid impact level "default | high | low", invalid type "fix | feature | chore", unknown section "<kebab-case of a module name> | <1st level dir in the repo>"
 - #3799 unknown section "monitoring"
 - #3809 invalid type "feat"

## Features


 - **[candi]** flow-schema new module [#3674](https://github.com/deckhouse/deckhouse/pull/3674)
    Added flow schema to prevent api overloading.
 - **[deckhouse]** Added bash wrapper for handling USR signals. [#3660](https://github.com/deckhouse/deckhouse/pull/3660)
 - **[deckhouse]** Added Python environment to support Python hooks [#3523](https://github.com/deckhouse/deckhouse/pull/3523)
 - **[deckhouse-config]** Support statuses for external modules [#3531](https://github.com/deckhouse/deckhouse/pull/3531)
    <what to expect for users, possibly MULTI-LINE>, required if impact_level is high ↓
 - **[istio]** Add istio version 1.16.2. [#3595](https://github.com/deckhouse/deckhouse/pull/3595)
    In environments where legacy versions of istio are used, the "D8IstioDeprecatedIstioVersionInstalled" alert will be fired.
 - **[runtime-audit-engine]** A new module to collect security events about possible threats in the cluster. [#3477](https://github.com/deckhouse/deckhouse/pull/3477)

## Fixes


 - **[candi]** Reorder swap disabling steps [#3772](https://github.com/deckhouse/deckhouse/pull/3772)
 - **[cloud-provider-vsphere]** Stop depending on CCM to uniquely identify instance ID. Fixes a couple of bugs. [#3721](https://github.com/deckhouse/deckhouse/pull/3721)
    <what to expect for users, possibly MULTI-LINE>, required if impact_level is high ↓
 - **[go_lib]** Remove go_lib/hooks/delete_not_matching_certificate_secret/hook.go [#3777](https://github.com/deckhouse/deckhouse/pull/3777)
 - **[log-shipper]** Bump librdkafka to v2.0.2 to make log-shipper read the full CA certificates chain for Kafka. [#3693](https://github.com/deckhouse/deckhouse/pull/3693)
 - **[log-shipper]** Make log-shipper-agents sending whole JSON message with metadata to Kafka destination. [#3692](https://github.com/deckhouse/deckhouse/pull/3692)
 - **[node-manager]** 1. Stop deleting Yandex.Cloud preemptible instances if percent of Ready Machines in a NodeGroup dips below 90%.
    2. Algorithm is simplified. [#3589](https://github.com/deckhouse/deckhouse/pull/3589)

## Chore


 - **[ceph-csi]** Added a small clarification to the migration docs [#3733](https://github.com/deckhouse/deckhouse/pull/3733)
 - **[docs]** Fix small typos in docs (DEVELOPMENT.md and 150-user-authn/docs/README_RU.md) [#3778](https://github.com/deckhouse/deckhouse/pull/3778)
 - **[log-shipper]** Update vector to `0.27.0`. [#3605](https://github.com/deckhouse/deckhouse/pull/3605)
 - **[monitoring-kubernetes]** DeprecatedDockerContainerRuntime alert is switched on — it is time to use containerd now. [#3763](https://github.com/deckhouse/deckhouse/pull/3763)

