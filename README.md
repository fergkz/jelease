# Jelease

![GitHub repo size](https://img.shields.io/github/repo-size/fergkz/jelease?style=for-the-badge)
![GitHub language count](https://img.shields.io/github/languages/count/fergkz/jelease?style=for-the-badge)
![GitHub forks](https://img.shields.io/github/forks/fergkz/jelease?style=for-the-badge)
![Bitbucket open issues](https://img.shields.io/bitbucket/issues/fergkz/jelease?style=for-the-badge)
![Bitbucket open pull requests](https://img.shields.io/bitbucket/pr-raw/fergkz/jelease?style=for-the-badge)


> Desenvolva um template (a princípio em twig) e extraia um html de release notes das suas sprints no jira

## 🏠 Ajustes e melhorias

O projeto ainda está em desenvolvimento

## 💻 Pré-requisitos

Antes de começar, verifique se você atendeu aos seguintes requisitos:

* Você instalou no mínimo a versão do `go 1.16`
* Você tem uma máquina `Windows / Linux / Mac`

## 🚀 Instalando o `Jelease`

Para instalar o `Jelease`, siga estas etapas:

* Instale o `go 1.16`
* Clone o repositório na sua máquina
* Instale os pacotes adicionais (ex: `go mod tidy`)
* Copie o arquivo `config-example.yml` e renomeie para `config.yml`
* Preencha as informações do arquivo `config.yml` com as suas configurações
* Copie o arquivo `template-example.twig` e renomeie para `template.twig`
* Edite o arquivo `template.twig` de acordo com o layout que deseja

## ☕ Compilando o `Jelease`

Você pode compilar este serviço para rodar diretamente na sua máquina, basta executar os comandos:

### Linux
```
> $Env:GOOS = "linux"; $Env:GOARCH = "amd64"
> go build -o launcher .
```

### Windows
```
> $Env:GOOS = "windows"; $Env:GOARCH = "amd64"
> go build -o launcher.exe .
```


## 🏁 Usando o `Jelease`

### Formatando mensagens
Para formatar mensagens de release notes, o `Jelease` está utilizando a seguinte formatação:

```
RELEASE NOTES

Title: <Título curto da nota de atualização>

Description: <Descrição amigável do que foi executado nesta nota de atualização>

Type: <Tipo de atualização: Ex: "Suporte / Sustentação", "Melhoria de usabilidade"...>

System: <Nome do sistema / módulo>
```
Basta fazer um comentário em uma atividade (atividade pai) no formato acima e já estará disponível no documento de atualização.

### Rodando o `Jelease` compilado:

* Basta executar o arquivo `laucher` no linux ou `launcher.exe` no windows
* Acesse do seu navegador a url `http://localhost/sprint/{NÚMERO_DA_SPRINT}`

### Rodando o `Jelease` em desenvolvimento:

* Execute o comando `go run .`
* Acesse do seu navegador a url `http://localhost/sprint/{NÚMERO_DA_SPRINT}`


## 😄 Seja uma das pessoas contribuidoras

Para contribuir com `Jelease`, siga estas etapas:

1. Bifurque este repositório.
2. Crie um branch: `git checkout -b main`.
3. Faça suas alterações e confirme-as: `git commit -m '<mensagem_commit>'`
4. Envie para o branch original: `git push origin jelease/main`
5. Crie a solicitação de pull.

Como alternativa, consulte a documentação do GitHub em [como criar uma solicitação pull](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request).


## 📝 Licença

Esse projeto está sob licença. Veja o arquivo [LICENÇA](LICENSE.md) para mais detalhes.

[⬆ Voltar ao topo](#jelease)<br>