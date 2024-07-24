# IAM 认证策略

- auto策略：

    该策略根据HTTP头`Authorization: Basic XX.ZZ.YY`和`Aut-horization: Bearer XX.YY.ZZ`自动选择Basic认证还是Bearer认证

- basic策略：

    实现Basic认证

- jwt策略：

    实现了Bearer认证，JWT是Bearer认证的具体实现

- cache策略：

    该策略是一个Bearer认证的实现，Token采用JWT格式，Token中的密钥ID存储在内存中，所以叫缓存认证

> [iam-authz](../../../../docs/iam/iam-authz.md)