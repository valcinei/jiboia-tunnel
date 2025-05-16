# 📦 Jiboia Tunnel — Estrutura e Uso

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

🌍 Esta documentação também está disponível em:
- [Español](README.es.md)
- [English](README.md)

O Jiboia Tunnel é uma ferramenta de tunelamento HTTP reverso baseada em WebSocket, inspirada em soluções como `ngrok` e `inlets`. A seguir está a documentação detalhada para desenvolvedores que desejam entender e replicar o projeto com precisão.

---

## 📁 Estrutura do Projeto (binários separados)

```
jiboia-tunnel/
├── cmd/
│   ├── jiboia-client/main.go     # Cliente: conecta ao relay e envia tráfego local
│   ├── jiboia-relay/main.go      # Relay: aceita WebSocket e encaminha requisições
│   ├── jiboia-server/main.go     # Mock: servidor local de teste
├── shared/
│   └── message.go                # Tipos comuns: TunnelMessage, TunnelResponse
├── go.mod
├── go.sum
└── README.md
```

Cada diretório `cmd/<nome>` define um binário separado:
- `jiboia-client`
- `jiboia-relay`
- `jiboia-server`

---

## 🚀 Comandos Disponíveis (pós-build)

### Iniciar o relay
```bash
./jiboia-relay
```
- WebSocket escutando em `/ws`
- Requisições HTTP são roteadas para os clientes conectados

---

### Iniciar o client (usuário local)
```bash
./jiboia-client http 3000
```
Atalho para expor `localhost:3000` como túnel via relay padrão (`ws://localhost:80/ws`).

Expõe seu serviço local (`localhost:3000`) como:
```
http://<nome-gerado>.jiboia.local
```

#### Com nome definido e relay remoto:
```bash
./jiboia-client http 3000 --name meuapp --relay wss://relay.jiboia.io/ws
```

#### Flags adicionais disponíveis:
| Flag             | Tipo     | Descrição                                                       |
|------------------|----------|-----------------------------------------------------------------|
| `--name`         | string   | Nome do túnel (subdomínio).                                     |
| `--relay`        | string   | Endereço WebSocket do relay.                                    |
| `--proto`        | string   | Protocolo a expor (`http`, `tcp`).                              |
| `--hostname`     | string   | Domínio customizado completo (ex: `meusite.com`).              |
| `--inspect`      | bool     | Mostra tráfego detalhado (modo debug).                          |
| `--authtoken`    | string   | Token de autenticação com o servidor.                          |
| `--config`       | string   | Caminho para arquivo de configuração externo.                   |
| `--region`       | string   | Região do relay (ex: `us`, `sa-east`).                         |
| `--label`        | string   | Identificador amigável do túnel (usado em logs/API futura).    |
| `--log-level`    | string   | Nível de log (`debug`, `info`, `warn`, `error`).               |

---

### Iniciar servidor local de teste
```bash
./jiboia-server
```
Responde com HTML simples em `http://localhost:3000`

---

## 🧪 Testando localmente com `go run`
```bash
# Terminal 1
sudo go run ./cmd/jiboia-relay/main.go

# Terminal 2
go run ./cmd/jiboia-server/main.go

# Terminal 3
go run ./cmd/jiboia-client/main.go --name jiboia --local http://localhost:3000
```

Abra no navegador:
```
http://jiboia.jiboia.local
```

Adicione ao seu `/etc/hosts`:
```
127.0.0.1 jiboia.jiboia.local
```

---

## 🛠 Buildando os binários
```bash
# Build todos manualmente
GOOS=linux GOARCH=amd64 go build -o jiboia-relay ./cmd/jiboia-relay
GOOS=linux GOARCH=amd64 go build -o jiboia-client ./cmd/jiboia-client
GOOS=linux GOARCH=amd64 go build -o jiboia-server ./cmd/jiboia-server
```
Ou com `goreleaser`, definindo múltiplos builds por binário.

---

## 🧱 Como funciona a aplicação
- **relay:** recebe requisições HTTP, extrai subdomínio, redireciona via WebSocket para um cliente conectado.
- **client:** escuta mensagens WebSocket e atua como proxy reverso para um servidor local.
- **server:** mock de aplicação para teste da cadeia de tunelamento.

---

## ✅ Etapas restantes para persistência e autenticação

### 🔐 Autenticação com JWT
1. Criar middleware `RequireAuth()` para proteger rotas (`/tunnels`, etc).
2. Aplicar o middleware às rotas REST no `jiboia-server`.
3. Adicionar validação do token JWT recebido via cookie ou `Authorization: Bearer`.
4. Criar endpoint opcional para `logout` (invalidar cookie).
5. (futuro) Criar rota `/users` com persistência de usuários.

### 💾 Persistência real em SQLite
1. Criar função `Migrate()` que execute `CREATE TABLE IF NOT EXISTS tunnels (...)`.
2. Criar tipo `SQLiteStore` que implemente a interface `TunnelStore`.
3. Substituir o uso de `InMemoryStore` por `SQLiteStore`.
4. Adicionar verificação de erro ao abrir o banco (permissão, caminho, etc).

### 🔑 Token no client
1. Adicionar suporte à flag `--authtoken` no `jiboia-client`.
2. Incluir token no header `Authorization: Bearer` ao fazer chamadas à API.
3. Validar token no `relay` para permitir ou negar conexão do túnel.

---

## 🌐 Suporte a domínios personalizados pelos usuários

### Objetivo
Permitir que o usuário, autenticado na plataforma, possa registrar e usar domínios personalizados para seus túneis.

### Etapas para suporte completo

7. **Cadastro do domínio pelo usuário (por usuário autenticado):**
   - Criar endpoint `POST /domains` na API do `jiboia-server`.
   - O domínio será vinculado ao túnel **e ao usuário autenticado**.
   - Evita conflito entre domínios e reforça segurança multiusuário.
   - Exemplo payload:
     ```json
     { "hostname": "meusite.com" }
     ```
   - Futuramente, exigir verificação via DNS TXT para validar propriedade.

8. **Validação no relay:**
   - O `relay` deve aceitar conexões por domínio customizado (não só subdomínios).
   - Verificar se o domínio existe no banco e está associado a um túnel ativo.

9. **Flags necessárias no client:**
   - `--hostname` para permitir domínios externos.

10. **Salvar no backend:**
    - Associar domínio ao túnel na base de dados.

11. **Configuração DNS:**
    - Usuário deve apontar o domínio para o IP do relay (registro A ou CNAME).

12. **HTTPS/TLS (futuro):**
    - Suporte com Let's Encrypt ou configuração manual com Nginx/Caddy.

13. **Validação DNS (futuro):**
    - Endpoint `/verify-domain` para confirmar propriedade via token.

---

## 🤝 Contribuindo com o projeto

Contribuições são bem-vindas! Você pode:
- Criar issues com ideias, bugs ou melhorias
- Enviar pull requests
- Participar das discussões

### Como começar
1. Faça um fork do repositório
2. Clone o fork localmente
3. Crie uma branch:
   ```bash
   git checkout -b minha-feature
   ```
4. Faça alterações e commite:
   ```bash
   git commit -m "feat: adiciona suporte a hostname personalizado"
   ```
5. Envie para seu fork:
   ```bash
   git push origin minha-feature
   ```
6. Abra um pull request para o repositório principal

Consulte [CONTRIBUTING.md](CONTRIBUTING.pt-BR.md)  para mais detalhes.

---

## 🔮 Futuras melhorias
- Autenticação com JWT
- HTTPS com Let's Encrypt
- Dashboard web de administração
- API REST no `jiboia-server` para gestão
- Balanceamento de carga entre relays

---

Essa separação por binários melhora a modularidade, facilita o deploy segmentado (ex: relay na nuvem e client local) e está pronta para produção ou expansão futura.
