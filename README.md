[![My Skills](https://skillicons.dev/icons?i=golang,docker,mysql,redis)](https://skillicons.dev)

## Balance service with mysql and redis
This app is created for testing purposes of production feature using go instead of one threaded lang.

### Main idea is to test updates in concurrent way and check if it's secure and correct
* avoid lost updates and dirty reads
* handle deadlocks
* secure transfer amount between balances
* secure deposit and withdraw
* add rate limit with redis
* redis queue if needed