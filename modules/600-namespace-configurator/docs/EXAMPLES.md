---
title: "The namespace-configurator module: examples"
---

## Example

This example will add `extended-monitoring.flant.com/enabled=true` label and `foo=bar` label to every namespace starting with `prod-` and `infra-`, except `infra-test`.

```yaml
namespaceConfiguratorEnabled: "true"
namespaceConfigurator: |
  configurations:
  - labels:
      extended-monitoring.flant.com/enabled: "true"
      foo: bar
    includeNames:
    - "prod-.*"
    - "infra-.*"
    excludeNames:
    - "infra-test"
```
