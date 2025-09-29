# ampctl

## What is it?

* ampctl is cli tool written in go to handle your amp (Apache MySQL PHP) stack
* ampctl using a system package manager as backend to install its packages.
* ampctl at moment only support **macs** with **homebrew** as package manager 

## PHP Extensions

We are using shivammathur/extensions tab, so check out the [repository](https://github.com/shivammathur/homebrew-extensions)


## Config

Example config under `~/.ampctl/config.yaml`

```
homebrew:
  path: /opt/homebrew
php:
  default: 8.3
  composer_versions: [2.2.25]
  versions:
    7.0:
      enabled: true
    7.1:
      enabled: true
    7.2:
      enabled: true
    7.3:
      enabled: true
    7.4:
      enabled: true
    8.0:
      enabled: true
    8.1:
      enabled: true
    8.2:
      enabled: true
    8.3:
      enabled: true
    8.4:
      enabled: true
apache:
  workspace: "/Users/acme/Workspace"
  http_port: 80
  https_port: 443
  ssl_certificate_cn: Acme Root R1
  ssl_certificate_county: DE
  ssl_certificate_locality: Munich
  ssl_certificate_organization: ACME
  ssl_certificate_organization_unit: dev
  ssl_certificate_province: Bavaria
database:
  versions:
    mariadb@10.11:
      enabled: true
      port: 3306
hosts:
  -
    host: dev.my-domain.com
    path: /Users/acme/Workspace/my-domain.com/public
    version: "8.3"
    ssl: true
```


