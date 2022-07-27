# О проекте

## Стек

* <img src="https://github.com/devicons/devicon/raw/master/icons/go/go-original.svg" width="15" height="15"> Go <img src="https://raw.githubusercontent.com/gin-gonic/logo/master/color.png" width="15" height="15"> [gin](https://github.com/gin-gonic/gin)
* <img src="https://github.com/devicons/devicon/raw/master/icons/postgresql/postgresql-original.svg" width="15" height="15"> PostgreSQL
* <img src="https://github.com/golangci/golangci-lint/raw/master/assets/go.png" width="15" height="15"> [golangci-lint](https://github.com/golangci/golangci-lint) <img src="https://camo.githubusercontent.com/c0bc16116647eb3c773360c495d8537d509df514fa8f77b545fca2edde5fc3d7/68747470733a2f2f6861646f6c696e742e6769746875622e696f2f6861646f6c696e742f696d672f6361745f636f6e7461696e65722e706e67" width="15" height="15"> <!--[hadolint](https://github.com/hadolint/hadolint) <img src="https://upload.wikimedia.org/wikipedia/commons/9/92/Yaml_logo.png" width="15" height="15"> yamllint-->
<!--* <img src="https://github.com/devicons/devicon/raw/master/icons/jenkins/jenkins-original.svg" width="15" height="15"> Jenkins -->
<!--* <img src="https://github.com/devicons/devicon/raw/master/icons/docker/docker-original.svg" width="15" height="15"> Docker -->
<!--* <img src="https://helm.sh/img/helm.svg" width="15" height="15"> Helm <img src="https://github.com/devicons/devicon/raw/master/icons/kubernetes/kubernetes-plain.svg" width="15" height="15"> Kubernetes <img src="https://raw.githubusercontent.com/kubernetes/minikube/master/images/logo/logo.png" width="15" height="15"> Minikube -->
<!--* <img src="https://github.com/devicons/devicon/raw/master/icons/prometheus/prometheus-original.svg" width="15" height="15"> Prometheus <img src="https://github.com/devicons/devicon/raw/master/icons/grafana/grafana-original.svg" width="15" height="15"> Grafana -->
<!-- * <img src="https://static-www.elastic.co/v3/assets/bltefdd0b53724fa2ce/blt36f2da8d650732a0/5d0823c3d8ff351753cbc99f/logo-elasticsearch-32-color.svg" width="15" height="15"> Elasticksearch <img src="https://static-www.elastic.co/v3/assets/bltefdd0b53724fa2ce/blt8b679e63f2b49b27/5d082d93877575d0584761c0/logo-logstash-32-color.svg" width="15" height="15"> Logstash <img src="https://static-www.elastic.co/v3/assets/bltefdd0b53724fa2ce/blt4466841eed0bf232/5d082a5e97f2babb5af907ee/logo-kibana-32-color.svg" width="15" height="15"> Kibana -->

# Задача

## 🛠 Сервис

1. Необходимо реализовать сервис, который должен предоставлять **API ленты постов**
2. Дистрибутивом сервиса должен быть <img src="https://github.com/devicons/devicon/raw/master/icons/docker/docker-original.svg" width="15" height="15"> Docker-образ
3. В качестве хранилища данных сервиса должна использоваться <img src="https://github.com/devicons/devicon/raw/master/icons/postgresql/postgresql-original.svg" width="15" height="15"> PostgreSQL <!-- ,события должны храниться в <img src="https://static-www.elastic.co/v3/assets/bltefdd0b53724fa2ce/blt36f2da8d650732a0/5d0823c3d8ff351753cbc99f/logo-elasticsearch-32-color.svg" width="15" height="15"> <img src="https://static-www.elastic.co/v3/assets/bltefdd0b53724fa2ce/blt8b679e63f2b49b27/5d082d93877575d0584761c0/logo-logstash-32-color.svg" width="15" height="15"> <img src="https://static-www.elastic.co/v3/assets/bltefdd0b53724fa2ce/blt4466841eed0bf232/5d082a5e97f2babb5af907ee/logo-kibana-32-color.svg" width="15" height="15"> ELK -->

### Endpoint'ы сервиса:

#### POST:

+ `/new-post` опубликовать новую запись

Пример тела запроса:
```json
{
	"content": "your text"
}
```

+ `/like?uuid=:uuid` лайкнуть запись с идентефикатором `uuid`

+ `/dislike?uuid=:uuid` дизлайкнуть запись с идентефикатором `uuid`

#### GET:

+ `/posts[?last=:number]` получить все записи / последние `:number`

Пример:
```json
{
	"total": 2,
	"data": [
		{
			"uuid": "1a",
			"content": "this is post a",
			"likes": 3,
			"dislikes": 2
		},
		{
			"uuid": "abacaba",
			"content": "abracadabra",
			"likes": 112,
			"dislikes": 0
		}
	]
}
```

+ `/healthz` получить статус о готовности сервиса

<!--
## ⚙️ CI/CD

* Конфигурационные файлы и исходный код должны пройти линтеры
* Код должен быть покрыт тестами
* Прохождение линтеров, юнит тестов и развертывание должно быть автоматизированно

## ☁️ Кластер

Сервис должен развертываться в <img src="https://github.com/devicons/devicon/raw/master/icons/kubernetes/kubernetes-plain.svg" width="15" height="15"> Kubernetes кластере
-->
<!--
## 📊 Мониторинг

В кластере должен быть настроен мониторинг сервиса и хранилища данных с помощью <img src="https://github.com/devicons/devicon/raw/master/icons/prometheus/prometheus-original.svg" width="15" height="15"> Prometheus, визуализация <img src="https://github.com/devicons/devicon/raw/master/icons/grafana/grafana-original.svg" width="15" height="15"> Grafana
* Health-checks
* Количество запросов в минуту
* Использование ресурсов кластера
-->
