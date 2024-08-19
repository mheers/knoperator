# knoperator Service

> Connects and combines nats mq with k8s

## Development

### Create credentials for MQ

```bash
nats-seeder --operator-seed SOAEUVFW77MXWK4IVK7RXXBMTWSVTFZFXKRJD5E26W622RUCO4CM7GZRHI --account-seed SAAPN4CYDRBUCWQ25HQLLWRO3ZS2HNCYWBZ7A2IF5AVGSTLZJJVP2EFMAU user-nkey -u test \
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
