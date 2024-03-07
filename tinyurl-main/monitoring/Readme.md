* Запуск

```
git clone https://github.com/skl256/grafana_stack_for_docker.git && \
cd grafana_stack_for_docker && \
sudo mkdir -p /mnt/common_volume/swarm/grafana/config && \
sudo mkdir -p /mnt/common_volume/grafana/{grafana-config,grafana-data,prometheus-data,loki-data,promtail-data} && \
sudo chown -R $(id -u):$(id -g) {/mnt/common_volume/swarm/grafana/config,/mnt/common_volume/grafana} && \
touch /mnt/common_volume/grafana/grafana-config/grafana.ini && \
cp config/* /mnt/common_volume/swarm/grafana/config/ && \
mv grafana.yaml docker-compose.yaml && \
docker compose up -d
```

* Перейти в браузере по адресу http:// **IP адрес сервера, на котором запущен стек** :3000 (логин `admin` пароль `admin`)
* Расширенные labels не будут работать до внесения изменений в `/etc/docker/daemon.json` (возможно, что файл не создан, соответственно необходимо создать самому)

* Необходимо добавить следуещее

{
   "log-driver": "json-file",
   "log-opts": {
     "labels-regex": "^.+"
  }
}
```
