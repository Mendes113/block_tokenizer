# Block Tokenizer

O projeto Block Tokenizer é voltado para startups júnior, proporcionando uma forma de democratizar o acesso ao mercado. A ideia central é converter startups em tokens, permitindo que qualquer pessoa possa comprar e monitorar o valor dessas startups.

## Estrutura do Projeto

O projeto é composto por três microserviços:
- 2 APIs em Go
- 1 gerador de blocos utilizando Rust
- 1 smart contract

## Como Executar?

Certifique-se de ter Go e Rust instalados em seu sistema.

### Kafka

Antes de executar os serviços, é necessário inicializar o servidor Kafka. Neste projeto, utilizamos o Conduktor para facilitar o setup, mas existem outras formas de executar o Kafka.

### API Primária

Para executar a API primária, siga os seguintes passos:

1. Clone o projeto:
    ```sh
    git clone <URL_DO_REPOSITORIO>
    ```

2. Acesse o diretório do projeto:
    ```sh
    cd <NOME_DO_DIRETORIO>
    ```

3. Acesse o diretório da API:
    ```sh
    cd api
    ```

4. Acesse o diretório `cmd`:
    ```sh
    cd cmd
    ```

5. Execute o comando para iniciar a API:
    ```sh
    go run main.go
    ```

### Configuração da API INFURA

Para que todos os endpoints funcionem corretamente, é necessário possuir uma API INFURA. Você pode se registrar e obter uma chave de API no site da INFURA [aqui](https://infura.io/).


