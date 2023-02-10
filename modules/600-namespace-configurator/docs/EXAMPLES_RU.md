---
title: "Модуль namespace-configurator: примеры"
---

## Пример

Этот пример добавит лейбл `extended-monitoring.flant.com/enabled=true` и label `foo=bar` к каждому Namespace, начинающемуся с `prod-` или `infra-`, за исключением `infra-test`.

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
