# Architecture Documentation

This document provides detailed architecture diagrams using Mermaid for the Passkey Demo application.

## Table of Contents
1. [System Architecture](#system-architecture)
2. [Registration Flow](#registration-flow)
3. [Authentication Flow](#authentication-flow)
4. [Component Interaction](#component-interaction)
5. [Deployment Architecture](#deployment-architecture)

---

## System Architecture

```mermaid
graph TB
    subgraph "Client Layer"
        Browser[Web Browser]
        WebAuthn[WebAuthn API]
        Authenticator[Authenticator<br/>Touch ID/Face ID/Security Key]
    end
    
    subgraph "Application Layer"
        Client[HTML/JavaScript Client]
        Server[Go Server]
        Handlers[HTTP Handlers]
        WebAuthnLib[go-webauthn Library]
    end
    
    subgraph "Data Layer"
        Storage[In-Memory Storage]
        Users[(Users)]
        Credentials[(Credentials)]
        Sessions[(Sessions)]
    end
    
    Browser --> Client
    Client --> WebAuthn
    WebAuthn --> Authenticator
    Client -->|HTTP/JSON| Server
    Server --> Handlers
    Handlers --> WebAuthnLib
    Handlers --> Storage
    Storage --> Users
    Storage --> Credentials
    Storage --> Sessions
    
    style Browser fill:#e1f5ff
    style Server fill:#fff4e1
    style Storage fill:#f0f0f0
```

---

## Registration Flow

```mermaid
sequenceDiagram
    participant User
    participant Browser
    participant Client as Client App
    participant Server as Go Server
    participant WebAuthn as WebAuthn Lib
    participant Storage
    participant Device as Authenticator

    User->>Browser: Enter username
    User->>Client: Click "Register"
    
    Client->>Server: POST /api/register/begin<br/>{username}
    Server->>Storage: Check/Create User
    Storage-->>Server: User object
    Server->>WebAuthn: BeginRegistration(user)
    WebAuthn-->>Server: Challenge + Options
    Server-->>Client: Credential Creation Options
    
    Client->>Browser: navigator.credentials.create()
    Browser->>Device: Request credential creation
    Device->>User: Prompt for biometric/PIN
    User->>Device: Provide authentication
    Device-->>Browser: New credential (public key)
    Browser-->>Client: Credential response
    
    Client->>Server: POST /api/register/finish<br/>{credential}
    Server->>WebAuthn: FinishRegistration(credential)
    WebAuthn->>WebAuthn: Verify attestation
    WebAuthn->>WebAuthn: Validate signature
    WebAuthn-->>Server: Validated credential
    Server->>Storage: Store credential
    Storage-->>Server: Success
    Server-->>Client: Registration successful
    Client-->>User: Show success message
```

---

## Authentication Flow

```mermaid
sequenceDiagram
    participant User
    participant Browser
    participant Client as Client App
    participant Server as Go Server
    participant WebAuthn as WebAuthn Lib
    participant Storage
    participant Device as Authenticator

    User->>Browser: Enter username
    User->>Client: Click "Login"
    
    Client->>Server: POST /api/login/begin<br/>{username}
    Server->>Storage: Get user + credentials
    Storage-->>Server: User with credentials
    Server->>WebAuthn: BeginLogin(user)
    WebAuthn-->>Server: Challenge + Options
    Server-->>Client: Credential Request Options
    
    Client->>Browser: navigator.credentials.get()
    Browser->>Device: Request authentication
    Device->>User: Prompt for biometric/PIN
    User->>Device: Provide authentication
    Device->>Device: Sign challenge with private key
    Device-->>Browser: Signed assertion
    Browser-->>Client: Assertion response
    
    Client->>Server: POST /api/login/finish<br/>{assertion}
    Server->>WebAuthn: FinishLogin(assertion)
    WebAuthn->>WebAuthn: Verify signature
    WebAuthn->>WebAuthn: Validate challenge
    WebAuthn->>WebAuthn: Check counter
    WebAuthn-->>Server: Validated credential
    Server->>Storage: Update credential (counter)
    Storage-->>Server: Success
    Server-->>Client: Login successful + user data
    Client-->>User: Show success + user info
```

---

## Component Interaction

```mermaid
graph LR
    subgraph "Frontend Components"
        Index[index.html<br/>Landing Page]
        Register[register.html<br/>Registration]
        Login[login.html<br/>Authentication]
    end
    
    subgraph "Backend Components"
        Main[main.go<br/>Server Setup]
        Handlers[handlers.go<br/>API Endpoints]
        Storage[storage.go<br/>Data Management]
    end
    
    subgraph "External Libraries"
        Mux[gorilla/mux<br/>Router]
        WebAuthnLib[go-webauthn<br/>WebAuthn Logic]
    end
    
    Index -->|Navigate| Register
    Index -->|Navigate| Login
    Register -->|API Call| Handlers
    Login -->|API Call| Handlers
    
    Main --> Mux
    Main --> WebAuthnLib
    Main --> Handlers
    Handlers --> Storage
    Handlers --> WebAuthnLib
    
    style Index fill:#e3f2fd
    style Register fill:#e3f2fd
    style Login fill:#e3f2fd
    style Main fill:#fff3e0
    style Handlers fill:#fff3e0
    style Storage fill:#fff3e0
```

---

## API Endpoint Flow

```mermaid
graph TD
    Start([Client Request]) --> Router{Route}
    
    Router -->|POST /api/register/begin| RegBegin[RegisterBegin Handler]
    Router -->|POST /api/register/finish| RegFinish[RegisterFinish Handler]
    Router -->|POST /api/login/begin| LoginBegin[LoginBegin Handler]
    Router -->|POST /api/login/finish| LoginFinish[LoginFinish Handler]
    Router -->|GET /health| Health[Health Check]
    Router -->|GET /*| Static[Static Files]
    
    RegBegin --> CreateUser{User Exists?}
    CreateUser -->|No| NewUser[Create User]
    CreateUser -->|Yes| ExistUser[Get User]
    NewUser --> GenChallenge1[Generate Challenge]
    ExistUser --> GenChallenge1
    GenChallenge1 --> SaveSession1[Save Session]
    SaveSession1 --> ReturnOptions1[Return Options]
    
    RegFinish --> GetSession1[Get Session]
    GetSession1 --> VerifyCred[Verify Credential]
    VerifyCred --> StoreCred[Store Credential]
    StoreCred --> CleanSession1[Delete Session]
    CleanSession1 --> ReturnSuccess1[Return Success]
    
    LoginBegin --> GetUser[Get User]
    GetUser --> GenChallenge2[Generate Challenge]
    GenChallenge2 --> SaveSession2[Save Session]
    SaveSession2 --> ReturnOptions2[Return Options]
    
    LoginFinish --> GetSession2[Get Session]
    GetSession2 --> VerifyAssertion[Verify Assertion]
    VerifyAssertion --> UpdateCred[Update Credential]
    UpdateCred --> CleanSession2[Delete Session]
    CleanSession2 --> ReturnSuccess2[Return Success + User]
    
    Health --> OK[Return OK]
    Static --> ServeFile[Serve HTML/CSS/JS]
    
    ReturnOptions1 --> End([Response])
    ReturnSuccess1 --> End
    ReturnOptions2 --> End
    ReturnSuccess2 --> End
    OK --> End
    ServeFile --> End
    
    style Start fill:#4caf50
    style End fill:#4caf50
    style Router fill:#2196f3
```

---

## Deployment Architecture

### Docker Deployment

```mermaid
graph TB
    subgraph "Docker Container"
        App[Go Application<br/>Port 8080]
        Client[Static Files<br/>/client]
    end
    
    subgraph "Host System"
        Docker[Docker Engine]
        Network[Bridge Network]
        Volume[Volume Mount<br/>Optional]
    end
    
    User[User Browser] -->|HTTP| Network
    Network -->|Port Mapping<br/>8080:8080| App
    App --> Client
    Docker --> App
    Volume -.->|Optional| App
    
    style App fill:#2196f3
    style Docker fill:#0db7ed
```

### Kubernetes Deployment

```mermaid
graph TB
    subgraph "Kubernetes Cluster"
        subgraph "Namespace: passkey-demo"
            Ingress[Ingress<br/>passkey.example.com]
            Service[Service<br/>passkey-service<br/>ClusterIP]
            
            subgraph "Deployment"
                Pod1[Pod 1<br/>passkey-app]
                Pod2[Pod 2<br/>passkey-app]
                Pod3[Pod 3<br/>passkey-app]
            end
            
            ConfigMap[ConfigMap<br/>app-config]
            Secret[Secret<br/>tls-cert]
        end
    end
    
    User[Users] -->|HTTPS| Ingress
    Ingress -->|TLS Termination| Service
    Service -->|Load Balance| Pod1
    Service --> Pod2
    Service --> Pod3
    
    ConfigMap -.->|Environment| Pod1
    ConfigMap -.->|Environment| Pod2
    ConfigMap -.->|Environment| Pod3
    
    Secret -.->|TLS Cert| Ingress
    
    style Ingress fill:#326ce5
    style Service fill:#326ce5
    style Pod1 fill:#4caf50
    style Pod2 fill:#4caf50
    style Pod3 fill:#4caf50
```

---

## Data Model

```mermaid
erDiagram
    USER ||--o{ CREDENTIAL : has
    USER ||--o{ SESSION : has
    
    USER {
        bytes ID
        string Name
        string DisplayName
    }
    
    CREDENTIAL {
        bytes ID
        bytes PublicKey
        bytes AttestationType
        bytes Transport
        bytes Flags
        bytes Authenticator
        uint32 SignCount
    }
    
    SESSION {
        string Username
        bytes Challenge
        string UserID
        timestamp Expires
    }
```

---

## Security Flow

```mermaid
graph TD
    Start([Request]) --> CORS{CORS Check}
    CORS -->|Valid| Auth{Endpoint Type}
    CORS -->|Invalid| Reject1[Reject: CORS]
    
    Auth -->|Public| Process[Process Request]
    Auth -->|Protected| ValidateSession{Session Valid?}
    
    ValidateSession -->|Yes| Process
    ValidateSession -->|No| Reject2[Reject: Unauthorized]
    
    Process --> Validate{Input Valid?}
    Validate -->|Yes| Execute[Execute Logic]
    Validate -->|No| Reject3[Reject: Bad Request]
    
    Execute --> WebAuthnCheck{WebAuthn Valid?}
    WebAuthnCheck -->|Yes| Success[Return Success]
    WebAuthnCheck -->|No| Reject4[Reject: Invalid Credential]
    
    Reject1 --> End([Error Response])
    Reject2 --> End
    Reject3 --> End
    Reject4 --> End
    Success --> End
    
    style Start fill:#4caf50
    style Success fill:#4caf50
    style End fill:#f44336
    style Reject1 fill:#f44336
    style Reject2 fill:#f44336
    style Reject3 fill:#f44336
    style Reject4 fill:#f44336
```

---

## Technology Stack

```mermaid
graph LR
    subgraph "Frontend"
        HTML[HTML5]
        CSS[CSS3]
        JS[JavaScript ES6+]
        WebAuthnAPI[WebAuthn API]
    end
    
    subgraph "Backend"
        Go[Go 1.21+]
        Mux[Gorilla Mux]
        WebAuthnGo[go-webauthn]
    end
    
    subgraph "Infrastructure"
        Docker[Docker]
        K8s[Kubernetes]
        Ingress[Ingress Controller]
    end
    
    subgraph "Security"
        FIDO2[FIDO2/WebAuthn]
        TLS[TLS/HTTPS]
        CORS[CORS]
    end
    
    HTML --> JS
    CSS --> JS
    JS --> WebAuthnAPI
    
    Go --> Mux
    Go --> WebAuthnGo
    
    Docker --> K8s
    K8s --> Ingress
    
    WebAuthnAPI --> FIDO2
    Ingress --> TLS
    Mux --> CORS
    
    style Go fill:#00add8
    style Docker fill:#0db7ed
    style K8s fill:#326ce5
    style FIDO2 fill:#ffab00
```

---

## Scalability Architecture

```mermaid
graph TB
    subgraph "Load Balancer Layer"
        LB[Load Balancer<br/>Nginx/HAProxy]
    end
    
    subgraph "Application Layer"
        App1[App Instance 1]
        App2[App Instance 2]
        App3[App Instance 3]
        AppN[App Instance N]
    end
    
    subgraph "Data Layer"
        Redis[(Redis<br/>Session Store)]
        DB[(PostgreSQL<br/>User/Credential Store)]
    end
    
    Users[Users] --> LB
    LB --> App1
    LB --> App2
    LB --> App3
    LB --> AppN
    
    App1 --> Redis
    App2 --> Redis
    App3 --> Redis
    AppN --> Redis
    
    App1 --> DB
    App2 --> DB
    App3 --> DB
    AppN --> DB
    
    style LB fill:#ff9800
    style App1 fill:#4caf50
    style App2 fill:#4caf50
    style App3 fill:#4caf50
    style AppN fill:#4caf50
    style Redis fill:#dc382d
    style DB fill:#336791
```

---

## Notes

- All diagrams are rendered using Mermaid syntax
- View these diagrams in any Markdown viewer that supports Mermaid
- GitHub, GitLab, and most modern documentation tools support Mermaid
- For best viewing experience, use a Mermaid-compatible viewer or IDE extension