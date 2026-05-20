# ⚙️ Attendify Backend

API REST du système de gestion des présences **Attendify**, développée en **Go** avec le framework **Gin**. Elle expose les endpoints consommés par l'application mobile employée et les interfaces d'administration.

---

## ✨ Fonctionnalités

### 🔐 Authentification (`/auth`)
| Méthode | Endpoint | Accès | Description |
|---|---|---|---|
| POST | `/auth/connexion-admin` | Public | Connexion administrateur |
| POST | `/auth/connexion-employe` | Public | Connexion employé |
| POST | `/auth/refresh-token` | Public | Renouvellement du token JWT |
| POST | `/auth/deconnexion` | Authentifié | Déconnexion |
| POST | `/auth/nouveau-compte-employe` | Admin | Créer un compte employé |
| PATCH | `/auth/reinitialiser-mot-de-passe/:id` | Admin | Réinitialiser le mot de passe d'un employé |

---

### 👤 Gestion des comptes (`/compte`)
| Méthode | Endpoint | Accès | Description |
|---|---|---|---|
| GET | `/compte/mes-informations` | Authentifié | Consulter son profil |
| PATCH | `/compte/modifier-son-mot-de-passe` | Authentifié | Changer son mot de passe |
| PATCH | `/compte/modifier-un-compte` | Authentifié | Modifier ses informations personnelles |
| GET | `/compte/tous-les-employes` | Admin | Lister tous les employés |
| GET | `/compte/tous-les-employes/:id` | Admin | Consulter le profil d'un employé |
| PATCH | `/compte/changer-de-site/:id/:siteID` | Admin | Affecter un employé à un site |
| DELETE | `/compte/supprimer-un-compte/:id` | Admin | Supprimer un compte employé |

---

### 🕐 Pointage (`/pointage`)
| Méthode | Endpoint | Accès | Description |
|---|---|---|---|
| POST | `/pointage/arrive` | Employé | Enregistrer une arrivée |
| PATCH | `/pointage/depart` | Employé | Enregistrer un départ |
| GET | `/pointage/tous-mes-pointages` | Employé | Consulter son historique de pointages |
| GET | `/pointage/est-sur-site` | Employé | Vérifier si l'employé est actuellement sur site |
| GET | `/pointage/tous-les-pointages` | Admin | Lister tous les pointages |
| GET | `/pointage/stats` | Admin | Consulter les statistiques de présence |
| GET | `/pointage/export` | Admin | Exporter les pointages (Excel) |
| DELETE | `/pointage/supprimer/:id` | Admin | Supprimer un pointage |

---

### 🌴 Congés (`/conge`)
| Méthode | Endpoint | Accès | Description |
|---|---|---|---|
| POST | `/conge/faire-une-demande` | Employé | Soumettre une demande de congé |
| GET | `/conge/tous-les-conges-employe` | Employé | Lister ses propres demandes |
| GET | `/conge/tous-les-conges-employe/:id` | Employé | Consulter le détail d'un congé |
| PATCH | `/conge/modifier/:id` | Employé | Modifier une demande de congé |
| GET | `/conge/tous-les-conges` | Admin | Lister toutes les demandes de congé |
| GET | `/conge/tous-les-conges/:id` | Admin | Consulter le détail d'un congé |
| PATCH | `/conge/modifier-statut/:id` | Admin | Approuver ou refuser un congé |
| DELETE | `/conge/supprimer/:id` | Admin | Supprimer une demande de congé |

---

### 📍 Sites (`/site`)
| Méthode | Endpoint | Accès | Description |
|---|---|---|---|
| GET | `/site/tous-les-sites/:id` | Authentifié | Consulter un site |
| POST | `/site/ajouter` | Admin | Créer un site |
| GET | `/site/tous-les-sites` | Admin | Lister tous les sites |
| PATCH | `/site/modifier/:id` | Admin | Modifier un site |
| DELETE | `/site/supprimer/:id` | Admin | Supprimer un site |

---

### 🔔 Notifications (`/notification`)
| Méthode | Endpoint | Accès | Description |
|---|---|---|---|
| GET | `/notification/toutes-les-notifications` | Authentifié | Lister ses notifications |
| PATCH | `/notification/modifier/:id` | Authentifié | Marquer une notification comme lue |
| DELETE | `/notification/supprimer/:id` | Authentifié | Supprimer une notification |

---

### 📲 Push Notifications (`/notification-push`)
| Méthode | Endpoint | Accès | Description |
|---|---|---|---|
| POST | `/notification-push` | Authentifié | Enregistrer ou mettre à jour un token push |
| DELETE | `/notification-push/:token` | Authentifié | Supprimer un token push |

---

## 🔑 Niveaux d'accès

| Niveau | Rôle | Description |
|---|---|---|
| Public | — | Aucune authentification requise |
| Authentifié | Tous | Token JWT valide requis |
| Employé | `role = 2` | Employé connecté |
| Admin | `role = 1` | Administrateur |

---

## 🛠️ Stack technique

| Technologie | Rôle |
|---|---|
| Go 1.25 | Langage principal |
| Gin | Framework HTTP |
| GORM | ORM (MySQL) |
| MySQL 8.0 | Base de données principale |
| Redis 8.4 | Cache / sessions |
| JWT (golang-jwt/jwt v5) | Authentification |
| Excelize | Export Excel |
| Docker + Docker Compose | Conteneurisation |

---

## 📁 Structure du projet

```
attendify-backend/
├── cmd/
│   └── api/
│       └── main.go        # Point d'entrée de l'application
├── configs/               # Fichiers de configuration (YAML)
├── internal/              # Logique métier (handlers, services, models)
├── pkg/
│   └── middlewares/       # AuthMiddleware, AutorisationMiddleware
├── Dockerfile             # Image multi-stage (builder + runtime Alpine)
├── docker-compose.yml     # MySQL, Redis, Adminer
├── go.mod / go.sum
└── .env
```

---

## 🚀 Démarrage rapide

### Prérequis

- [Go 1.25+](https://go.dev/dl/)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)

### 1. Cloner le dépôt

```bash
git clone https://github.com/cabitibaly/attendify-backend.git
cd attendify-backend
```

### 2. Configurer les variables d'environnement

Créez un fichier `.env` à la racine :

```env
# Serveur
PORT=8080

# Base de données
DB_HOST=localhost
DB_PORT=3312
DB_NAME=attendify
DB_USER=root
DB_PASSWORD=root

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your_secret_key
```

### 3. Lancer les services (MySQL + Redis)

```bash
docker compose up -d
```

| Service | URL |
|---|---|
| MySQL | `localhost:3312` |
| Redis | `localhost:6379` |
| Adminer (UI DB) | [http://localhost:7080](http://localhost:7080) |

### 4. Lancer le serveur

```bash
go mod download
go run ./cmd/api/main.go
```

L'API est accessible sur `http://localhost:8080`.

---

## 🐳 Déploiement Docker

```bash
docker build -t attendify-backend .
docker run -p 8080:8080 --env-file .env attendify-backend
```

---

## 🧪 Tests

```bash
go test ./...
```

---

## Auteur
 
**cabitibaly** — [GitHub](https://github.com/cabitibaly)

## Liens vers les apps mobiles

[Attendify-employe](https://github.com/cabitibaly/attendify-employee)
[Attendify-admin](https://github.com/cabitibaly/attendify-admin)
