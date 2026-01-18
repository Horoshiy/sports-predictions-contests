# Architecture Diagram

## System Architecture

```mermaid
graph TB
    %% External Users
    User[👤 Users]
    TelegramUser[👤 Telegram Users]
    
    %% Frontend Layer
    Frontend[🌐 React Frontend<br/>Port 3000]
    TelegramBot[🤖 Telegram Bot]
    
    %% API Gateway Layer
    APIGateway[🚪 API Gateway<br/>Port 8080<br/>HTTP REST]
    
    %% Microservices Layer
    UserService[🔐 User Service<br/>Port 8084<br/>Authentication & Users]
    ContestService[🏆 Contest Service<br/>Port 8085<br/>Contests & Teams]
    PredictionService[🎯 Prediction Service<br/>Port 8086<br/>Predictions & Events]
    ScoringService[📊 Scoring Service<br/>Port 8087<br/>Scores & Leaderboards]
    SportsService[⚽ Sports Service<br/>Port 8088<br/>Sports Data & Sync]
    NotificationService[📢 Notification Service<br/>Port 8089<br/>Multi-channel Notifications]
    
    %% Data Layer
    PostgreSQL[(🗄️ PostgreSQL<br/>Port 5432<br/>Primary Database)]
    Redis[(⚡ Redis<br/>Port 6379<br/>Cache & Sessions)]
    
    %% External Services
    TheSportsDB[🌐 TheSportsDB API<br/>External Sports Data]
    TelegramAPI[📱 Telegram Bot API]
    EmailSMTP[📧 SMTP Server<br/>Email Notifications]
    
    %% User Connections
    User --> Frontend
    User --> APIGateway
    TelegramUser --> TelegramBot
    
    %% Frontend Connections
    Frontend --> APIGateway
    TelegramBot --> UserService
    TelegramBot --> ContestService
    TelegramBot --> ScoringService
    TelegramBot --> NotificationService
    
    %% API Gateway Routing
    APIGateway --> UserService
    APIGateway --> ContestService
    APIGateway --> PredictionService
    APIGateway --> ScoringService
    APIGateway --> SportsService
    APIGateway --> NotificationService
    
    %% Inter-service Communication (gRPC)
    PredictionService -.->|gRPC| ContestService
    ScoringService -.->|gRPC| ContestService
    ScoringService -.->|gRPC| PredictionService
    NotificationService -.->|gRPC| UserService
    
    %% Database Connections
    UserService --> PostgreSQL
    ContestService --> PostgreSQL
    PredictionService --> PostgreSQL
    ScoringService --> PostgreSQL
    SportsService --> PostgreSQL
    NotificationService --> PostgreSQL
    
    %% Cache Connections
    UserService --> Redis
    ContestService --> Redis
    ScoringService --> Redis
    NotificationService --> Redis
    
    %% External API Connections
    SportsService --> TheSportsDB
    TelegramBot --> TelegramAPI
    NotificationService --> EmailSMTP
    
    %% Styling
    classDef frontend fill:#e1f5fe
    classDef gateway fill:#f3e5f5
    classDef service fill:#e8f5e8
    classDef database fill:#fff3e0
    classDef external fill:#fce4ec
    
    class Frontend,TelegramBot frontend
    class APIGateway gateway
    class UserService,ContestService,PredictionService,ScoringService,SportsService,NotificationService service
    class PostgreSQL,Redis database
    class TheSportsDB,TelegramAPI,EmailSMTP external
```

## Service Communication Flow

```mermaid
sequenceDiagram
    participant U as User
    participant F as Frontend
    participant AG as API Gateway
    participant US as User Service
    participant CS as Contest Service
    participant PS as Prediction Service
    participant SS as Scoring Service
    participant NS as Notification Service
    participant DB as PostgreSQL
    participant R as Redis
    
    %% User Registration Flow
    U->>F: Register Account
    F->>AG: POST /v1/auth/register
    AG->>US: Forward Request
    US->>DB: Store User Data
    US->>AG: Return JWT Token
    AG->>F: Return Response
    F->>U: Show Success
    
    %% Contest Creation Flow
    U->>F: Create Contest
    F->>AG: POST /v1/contests (with JWT)
    AG->>CS: Forward Request
    CS->>DB: Store Contest
    CS->>AG: Return Contest Data
    AG->>F: Return Response
    F->>U: Show Contest Created
    
    %% Prediction Submission Flow
    U->>F: Submit Prediction
    F->>AG: POST /v1/predictions (with JWT)
    AG->>PS: Forward Request
    PS->>CS: Validate Contest Access
    CS->>PS: Return Validation
    PS->>DB: Store Prediction
    PS->>AG: Return Prediction Data
    AG->>F: Return Response
    F->>U: Show Prediction Submitted
    
    %% Scoring Flow
    Note over SS: Event Result Available
    SS->>PS: Get Predictions for Event
    PS->>SS: Return Predictions
    SS->>DB: Calculate and Store Scores
    SS->>R: Update Leaderboard Cache
    SS->>NS: Send Score Notifications
    NS->>U: Deliver Notifications
```

## Data Flow Architecture

```mermaid
flowchart LR
    %% Data Sources
    ExtAPI[External Sports APIs]
    UserInput[User Input]
    
    %% Processing Layer
    SportsSync[Sports Data Sync]
    PredictionEngine[Prediction Engine]
    ScoringEngine[Scoring Engine]
    
    %% Storage Layer
    PrimaryDB[(Primary Database)]
    CacheLayer[(Cache Layer)]
    
    %% Output Layer
    WebUI[Web Interface]
    MobileAPI[Mobile API]
    TelegramBot[Telegram Bot]
    Notifications[Notifications]
    
    %% Data Flow
    ExtAPI --> SportsSync
    SportsSync --> PrimaryDB
    
    UserInput --> PredictionEngine
    PredictionEngine --> PrimaryDB
    
    PrimaryDB --> ScoringEngine
    ScoringEngine --> PrimaryDB
    ScoringEngine --> CacheLayer
    
    PrimaryDB --> WebUI
    CacheLayer --> WebUI
    PrimaryDB --> MobileAPI
    CacheLayer --> MobileAPI
    PrimaryDB --> TelegramBot
    
    ScoringEngine --> Notifications
    Notifications --> WebUI
    Notifications --> TelegramBot
```

## Network Architecture

```mermaid
graph TB
    %% Load Balancer
    LB[🔄 Load Balancer<br/>nginx/traefik]
    
    %% Application Tier
    subgraph "Application Tier"
        AG1[API Gateway 1]
        AG2[API Gateway 2]
        
        subgraph "Microservices Cluster"
            US[User Service]
            CS[Contest Service]
            PS[Prediction Service]
            SS[Scoring Service]
            SpS[Sports Service]
            NS[Notification Service]
        end
    end
    
    %% Data Tier
    subgraph "Data Tier"
        DBMaster[(PostgreSQL Master)]
        DBReplica[(PostgreSQL Replica)]
        RedisCluster[(Redis Cluster)]
    end
    
    %% External Services
    subgraph "External Services"
        CDN[CDN]
        Monitoring[Monitoring]
        Logging[Centralized Logging]
    end
    
    %% Connections
    Internet --> LB
    LB --> AG1
    LB --> AG2
    
    AG1 --> US
    AG1 --> CS
    AG1 --> PS
    AG2 --> SS
    AG2 --> SpS
    AG2 --> NS
    
    US --> DBMaster
    CS --> DBMaster
    PS --> DBReplica
    SS --> DBMaster
    SpS --> DBMaster
    NS --> DBMaster
    
    US --> RedisCluster
    CS --> RedisCluster
    SS --> RedisCluster
    NS --> RedisCluster
    
    LB --> CDN
    AG1 --> Monitoring
    AG2 --> Monitoring
    US --> Logging
    CS --> Logging
```

## Deployment Architecture

```mermaid
graph TB
    %% Development Environment
    subgraph "Development"
        DevDocker[Docker Compose<br/>All Services Local]
        DevDB[(Local PostgreSQL)]
        DevRedis[(Local Redis)]
    end
    
    %% Staging Environment
    subgraph "Staging"
        StagingK8s[Kubernetes Cluster]
        StagingDB[(Managed Database)]
        StagingRedis[(Managed Redis)]
    end
    
    %% Production Environment
    subgraph "Production"
        ProdK8s[Kubernetes Cluster<br/>Multi-AZ]
        ProdDB[(RDS/CloudSQL<br/>Multi-AZ)]
        ProdRedis[(ElastiCache/MemoryStore)]
        ProdCDN[CloudFront/CloudFlare]
        ProdMonitoring[Prometheus/Grafana]
    end
    
    %% CI/CD Pipeline
    Git[Git Repository] --> CI[GitHub Actions]
    CI --> DevDocker
    CI --> StagingK8s
    CI --> ProdK8s
    
    %% Data Flow
    DevDocker --> DevDB
    DevDocker --> DevRedis
    
    StagingK8s --> StagingDB
    StagingK8s --> StagingRedis
    
    ProdK8s --> ProdDB
    ProdK8s --> ProdRedis
    ProdK8s --> ProdCDN
    ProdK8s --> ProdMonitoring
```

## Security Architecture

```mermaid
graph TB
    %% External Layer
    Internet[Internet]
    WAF[Web Application Firewall]
    
    %% Security Layer
    subgraph "Security Layer"
        SSL[SSL/TLS Termination]
        RateLimit[Rate Limiting]
        Auth[JWT Authentication]
    end
    
    %% Application Layer
    subgraph "Application Layer"
        APIGateway[API Gateway]
        Services[Microservices]
    end
    
    %% Data Layer
    subgraph "Data Layer"
        EncryptedDB[(Encrypted Database)]
        SecureRedis[(Secured Redis)]
    end
    
    %% Security Controls
    subgraph "Security Controls"
        Secrets[Secrets Management]
        Audit[Audit Logging]
        Monitoring[Security Monitoring]
    end
    
    %% Flow
    Internet --> WAF
    WAF --> SSL
    SSL --> RateLimit
    RateLimit --> Auth
    Auth --> APIGateway
    APIGateway --> Services
    Services --> EncryptedDB
    Services --> SecureRedis
    
    Services --> Secrets
    Services --> Audit
    Services --> Monitoring
```

## Technology Stack Overview

| Layer | Technology | Purpose |
|-------|------------|---------|
| **Frontend** | React, TypeScript, Material-UI | User interface |
| **API Gateway** | Go, gRPC-Gateway | HTTP/gRPC translation |
| **Microservices** | Go, gRPC | Business logic |
| **Database** | PostgreSQL 15 | Primary data storage |
| **Cache** | Redis 7 | Session & leaderboard cache |
| **Message Queue** | Redis Streams | Async processing |
| **Monitoring** | Prometheus, Grafana | Metrics & alerting |
| **Logging** | Structured JSON logs | Centralized logging |
| **Deployment** | Docker, Kubernetes | Container orchestration |
| **CI/CD** | GitHub Actions | Automated deployment |
| **Security** | JWT, TLS, WAF | Authentication & protection |

---

*These diagrams provide a visual overview of the Sports Prediction Contests platform architecture. For implementation details, refer to the service-specific documentation.*
