# 📦 Jiboia Tunnel — Estructura y Uso

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

🌍 Esta documentación también está disponible en:
- [Português (Brasil)](README.pt-BR.md)
- [English](README.md)

Jiboia Tunnel es una herramienta de túnel HTTP inverso basada en WebSocket, inspirada en soluciones como `ngrok` e `inlets`. A continuación se presenta una documentación detallada para desarrolladores que deseen comprender y replicar el proyecto con precisión.

---

## 📁 Estructura del Proyecto (binarios separados)

```
jiboia-tunnel/
├── cmd/
│   ├── jiboia-client/main.go     # Cliente: se conecta al relay y envía tráfico local
│   ├── jiboia-relay/main.go      # Relay: acepta WebSocket y reenvía solicitudes
│   ├── jiboia-server/main.go     # Mock: servidor local de prueba
├── shared/
│   └── message.go                # Tipos comunes: TunnelMessage, TunnelResponse
├── go.mod
├── go.sum
└── README.md
```

Cada directorio `cmd/<nombre>` define un binario separado:
- `jiboia-client`
- `jiboia-relay`
- `jiboia-server`

---

## 🚀 Comandos Disponibles (post-build)

### Iniciar el relay
```bash
./jiboia-relay
```
- WebSocket escuchando en `/ws`
- Las solicitudes HTTP se enrutan a los clientes conectados

---

### Iniciar el cliente (usuario local)
```bash
./jiboia-client http 3000
```
Acceso directo para exponer `localhost:3000` como túnel mediante relay por defecto (`ws://localhost:80/ws`).

Expone su servicio local (`localhost:3000`) como:
```
http://<nombre-generado>.jiboia.local
```

#### Con nombre definido y relay remoto:
```bash
./jiboia-client http 3000 --name miapp --relay wss://relay.jiboia.io/ws
```

#### Flags adicionales disponibles:
| Flag             | Tipo     | Descripción                                                    |
|------------------|----------|----------------------------------------------------------------|
| `--name`         | string   | Nombre del túnel (subdominio).                                |
| `--relay`        | string   | Dirección WebSocket del relay.                                |
| `--proto`        | string   | Protocolo a exponer (`http`, `tcp`).                          |
| `--hostname`     | string   | Dominio personalizado completo (ej: `misitio.com`).           |
| `--inspect`      | bool     | Modo detallado de tráfico (debug).                            |
| `--authtoken`    | string   | Token de autenticación con el servidor.                      |
| `--config`       | string   | Ruta a archivo de configuración externo.                      |
| `--region`       | string   | Región del relay (ej: `us`, `sa-east`).                       |
| `--label`        | string   | Etiqueta identificadora del túnel (visible en logs/API).      |
| `--log-level`    | string   | Nivel de log (`debug`, `info`, `warn`, `error`).              |

---

### Iniciar servidor local de prueba
```bash
./jiboia-server
```
Responde con HTML simple en `http://localhost:3000`

---

## 🧪 Pruebas locales con `go run`
```bash
# Terminal 1
sudo go run ./cmd/jiboia-relay/main.go

# Terminal 2
go run ./cmd/jiboia-server/main.go

# Terminal 3
go run ./cmd/jiboia-client/main.go --name jiboia --local http://localhost:3000
```

Abra en el navegador:
```
http://jiboia.jiboia.local
```

Agregue en su `/etc/hosts`:
```
127.0.0.1 jiboia.jiboia.local
```

---

## 🛠 Compilación de binarios
```bash
# Compilar todos manualmente
GOOS=linux GOARCH=amd64 go build -o jiboia-relay ./cmd/jiboia-relay
GOOS=linux GOARCH=amd64 go build -o jiboia-client ./cmd/jiboia-client
GOOS=linux GOARCH=amd64 go build -o jiboia-server ./cmd/jiboia-server
```
O con `goreleaser`, definiendo múltiples builds por binario.

---

## 🧱 ¿Cómo funciona la aplicación?
- **relay:** recibe solicitudes HTTP, extrae el subdominio, redirige vía WebSocket al cliente conectado.
- **client:** escucha mensajes WebSocket y actúa como proxy reverso para un servidor local.
- **server:** aplicación mock para pruebas de tunelamiento.

---

## ✅ Pasos restantes para persistencia y autenticación

### 🔐 Autenticación JWT
1. Crear middleware `RequireAuth()` para proteger rutas (`/tunnels`, etc).
2. Aplicar el middleware a las rutas REST en `jiboia-server`.
3. Validar el token JWT recibido por cookie o `Authorization: Bearer`.
4. Crear endpoint opcional para `logout`.
5. (futuro) Crear ruta `/users` con persistencia de usuarios.

### 💾 Persistencia real con SQLite
1. Crear función `Migrate()` que ejecute `CREATE TABLE IF NOT EXISTS tunnels (...)`.
2. Crear tipo `SQLiteStore` que implemente la interfaz `TunnelStore`.
3. Reemplazar `InMemoryStore` por `SQLiteStore`.
4. Verificar errores al abrir la base de datos (ruta, permisos).

### 🔑 Token en el cliente
1. Agregar flag `--authtoken` en `jiboia-client`.
2. Incluir token en el header `Authorization: Bearer` al llamar la API.
3. Validar token en el relay para permitir o denegar conexión.

---

## 🌐 Soporte para dominios personalizados

### Objetivo
Permitir que el usuario autenticado registre y utilice dominios propios.

### Pasos para implementación

7. **Registro del dominio por el usuario:**
   - Crear endpoint `POST /domains` en `jiboia-server`.
   - Asociar el dominio al túnel y al usuario autenticado.
   - Validar dominio vía DNS TXT (opcional).

8. **Validación en el relay:**
   - Aceptar `Host: misitio.com`, además de subdominios.
   - Confirmar existencia del dominio en base de datos.

9. **Flags necesarias en client:**
   - `--hostname` para usar dominio externo.

10. **Guardar en backend:**
    - Asociar dominio al túnel en la base de datos.

11. **Configuración DNS:**
    - Usuario debe apuntar su dominio a la IP del relay (registro A o CNAME).

12. **HTTPS/TLS (futuro):**
    - Soporte con Let's Encrypt o configuración manual con Nginx/Caddy.

13. **Validación DNS (futuro):**
    - Endpoint `/verify-domain` para confirmar propiedad vía token.

---

## 🤝 Contribuyendo al proyecto

¡Contribuciones bienvenidas! Puedes:
- Crear issues con ideas, errores o mejoras
- Enviar pull requests
- Participar en debates

### Cómo empezar
1. Haz un fork del repositorio
2. Clona el fork
3. Crea una rama:
   ```bash
   git checkout -b mi-feature
   ```
4. Realiza los cambios y haz commit:
   ```bash
   git commit -m "feat: soporte para hostname personalizado"
   ```
5. Haz push:
   ```bash
   git push origin mi-feature
   ```
6. Abre un pull request al repositorio principal

Consulta [CONTRIBUTING.md](CONTRIBUTING.es.md)  para más detalles.

---

## 🔮 Futuras mejoras
- Autenticación con JWT
- HTTPS con Let's Encrypt
- Dashboard web de administración
- API REST para gestión
- Balanceo de carga entre relays

---

Esta separación por binarios mejora la modularidad, facilita el despliegue distribuido (ej: relay en la nube y cliente local) y está lista para entornos de producción o expansión futura.
