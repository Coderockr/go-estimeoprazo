# Estime o Prazo

Serviço que faz simulações usando o método de Monte Carlo para encontrar a probabilidade de finalização de um projeto

### Testes unitários

Para rodar os testes unitários é necessário:

```sh
cd estimeoprazo
export GOPATH=`pwd`
export GO111MODULE=off
go test
```

### Execução local

É necessário instalar o Go e o [SDK do Google App Engine](https://cloud.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go), devido ao projeto estar configurado para rodar neste cloud server

```sh
cd estimeoprazo
export GOPATH=`pwd`
export GO111MODULE=off
goapp serve
```

O serviço vai ouvir no endereço `http://localhost:8080`

### Execução no Google App Engine

Para fazer deploy de uma nova versão do app é necessário:

```sh
git add .
git commit -m "Alterações"
git push origin master
goapp deploy --application $YOUR_GOOGLE_APP estimeoprazo
```

O serviço vai ouvir no endereço `https://$YOUR_GOOGLE_APP.appspot.com`

O serviço espera requisições via `POST` com o `header`

    Content-Type: application/x-www-form-urlencoded

E no `body` os parâmetros, como no exemplo abaixo:

    MinTasks=200&MaxTasks=220&MinSplitTasks=1&MaxSplitTasks=3&MinTasksDone=30&MaxTasksDone=40
