# knoperator Service

> Connects and combines nats mq with k8s

## Development

### Create credentials for MQ

```bash
nats-seeder --operator-seed SOAA7B6L7CCSKP7XFDD6MXH65SRZIHBL2HXZFHWOLSVTN3UQ4EBMHKAJ2A --account-seed SAAGQMOVIAG5TTHENP3JMN3HOA4YY3LXJSU6SAWUDTPSIECKKSC54QYOJQ user-nkey -u test \
-p "knoperator.pods.create" \
-p "knoperator.pods.get" \
-p "knoperator.pods.update" \
-p "knoperator.pods.delete" \
-p "knoperator.pods.watch" \
-p "knoperator.deployments.create" \
-p "knoperator.deployments.get" \
-p "knoperator.deployments.update" \
-p "knoperator.deployments.delete" \
-p "knoperator.deployments.watch" \
-p "knoperator.deployments.scale" \
-s "knoperator.*.*" \
-s "knoperator.*.*.*" \
-p "_INBOX.>" \
-s "_INBOX.>" > test.creds
```
