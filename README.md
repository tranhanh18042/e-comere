

## run system

    docker compose up --build

## run load test

    cd load_test

    ddosify --config <file_config>

## import dashboard Grafana


go to the link: http://localhost:9981/

default account: admin/ admin

import file .json : /infra/prometheus/dashboards/Mornitoring infrastructure.json
import file .json : /infra/prometheus/dashboards/Mornitoring system ecomere.json