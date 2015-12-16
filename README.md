# configlog

 Package for simpe config load and enable logfile
 If logfile not found, log output will be printed to STDOUT

###chanelog

####16.12.2015
  
  For support logrotate added log file watcher. It watch for rename, delete.  

####01.12.2015
  
  Disabled reload by time.
  
  
####use:

```
$ go get github.com/andboson/configlog
```

sample dir structure:

```
config/app.yml
      /production/app.yml
tests/default_tests.go
main.go
```

If folder  `production` exists config file will be loaded form there

Sample config file

```
port: 8001
debug: true
logfile: 'log/logfile.log'
database:
    host: '127.0.0.1'
    pg_base: pgbase
    pg_user: root
    pg_pass: root
    pg_host: '127.0.0.1'
    db_port: 5432
    sslmode: disable
    sslcert: ''
    sslkey: ''
    sslrootcert: ''

redis:
    redis_host: '127.0.0.1'
    redis_port: 6379
    redis_pass: ''
```
log file loaded from `logfile` config value

Sample use:
```
import (
    . "github.com/andboson/configlog"
)



 debug, error := AppConfig.String("debug")
```

more see <a href="github.com/olebedev/config">github.com/olebedev/config</a>