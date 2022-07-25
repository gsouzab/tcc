# Trabalho de Conclusão de Curso - DCC/UFRJ 2018/1

Alunos:
- Emerson Yamamoto dos Santos
- Gabriel de Souza Bustamante

Pré requisitos:
- [Docker](https://docs.docker.com/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Yarn](https://yarnpkg.com/en/docs/install) ou [npm](https://www.npmjs.com/get-npm)

Para executar o projeto:
``` sh
git clone https://github.com/gsouzab/tcc  #clonando repositório
cd <DIR_RAIZ_DO_PROJETO>/client && yarn #instalando dependências (pode substituir yarn por npm install)
cd <DIR_RAIZ_DO_PROJETO> && docker compose up -d #iniciando containers
```

<sup> Onde <DIR_RAIZ_DO_PROJETO> é o caminho da raiz do projeto ex: **/home/gabriel-bustamante/workspace/tcc** </sup>

Para parar:
``` sh
cd <DIR_RAIZ_DO_PROJETO> && docker compose down #parando containers
```

A aplicação cliente deve ficar disponível em [localhost:8080](http://localhost:8080) e a interface de gerenciamento dos containers Docker (utilizando [portainer](https://github.com/portainer/portainer)) em [localhost:10001](http://localhost:10001)
