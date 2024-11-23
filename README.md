# knoperator Service

> Connects and combines nats mq with k8s

## Development

### Create credentials for MQ

```bash
nats-seeder --operator-seed SOAEUVFW77MXWK4IVK7RXXBMTWSVTFZFXKRJD5E26W622RUCO4CM7GZRHI --account-seed SAAPN4CYDRBUCWQ25HQLLWRO3ZS2HNCYWBZ7A2IF5AVGSTLZJJVP2EFMAU user-nkey -u knoperator \
-p "knoperator.*.*" \
-p "knoperator.*.*.*" \
-s "knoperator.*.*" \
-s "knoperator.*.*.*" \
-p "_INBOX.>" \
-s "_INBOX.>" > mq.creds
```

# TODO
- [ ] write github actions pipeline using dagger
