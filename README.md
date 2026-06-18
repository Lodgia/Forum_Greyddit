# Forum_Greyddit

Forum communautaire inspiré de Reddit, développé en Go avec SQLite.

---

## Sommaire

- [Forum\_Greyddit](#forum_greyddit)
  - [Sommaire](#sommaire)
  - [Présentation](#présentation)
  - [Fonctionnalités](#fonctionnalités)
    - [Sans être connecté](#sans-être-connecté)
    - [En étant connecté](#en-étant-connecté)
  - [Structure du projet](#structure-du-projet)
  - [Installation](#installation)
    - [Prérequis](#prérequis)
    - [Cloner le projet](#cloner-le-projet)
    - [Installer les dépendances](#installer-les-dépendances)
    - [Lancer le serveur](#lancer-le-serveur)
  - [Utilisation](#utilisation)
    - [Accéder au site](#accéder-au-site)
  - [Base de données](#base-de-données)
  - [Sécurité](#sécurité)
  - [Technologies utilisées](#technologies-utilisées)
  - [Auteurs](#auteurs)

---

## Présentation

Greyddit est une application web de forum permettant aux utilisateurs de créer des posts, de les catégoriser, de commenter et de voter. Les visiteurs non connectés peuvent lire le contenu. Les utilisateurs inscrits peuvent interagir pleinement avec la communauté.

---

## Fonctionnalités

### Sans être connecté
- Lire tous les posts et commentaires
- Voir les scores de likes / dislikes
- Naviguer par catégorie
- Accéder aux pages Connexion et Inscription

### En étant connecté
- Créer des posts avec une ou plusieurs catégories
- Modifier et supprimer ses propres posts
- Commenter les posts
- Supprimer ses propres commentaires
- Liker ou disliker les posts et commentaires
- Filtrer les posts : **Tous**, **Mes posts**, **Likés**

---

## Structure du projet

```
Forum_Greyddit/
├── main.go                  # Point d'entrée, routing, démarrage du serveur
├── go.mod                   # Dépendances Go
├── go.sum                   # Checksums des dépendances
├── forum.db                 # Base de données SQLite (généré au premier lancement)
│
├── database/
│   ├── db.go                # Connexion à SQLite, initialisation
│   ├── schema.sql           # Création des tables et données par défaut
│   └── queries.go           # Toutes les requêtes SQL (CRUD)
│
├── handlers/
│   ├── auth.go              # Inscription, connexion, déconnexion, rendu des templates
│   ├── posts.go             # Affichage, création, modification, suppression des posts
│   └── comments.go          # Commentaires et votes (posts + commentaires)
│
├── middleware/
│   └── session.go           # Lecture du cookie, injection utilisateur, RequireAuth
│
├── models/
│   └── models.go            # Structures de données : User, Post, Comment, Category, Session
│
├── templates/
│   ├── base.html            # Squelette commun, navbar
│   ├── index.html           # Page d'accueil, liste des posts
│   ├── post.html            # Page de détail d'un post et ses commentaires
│   ├── create.html          # Formulaire de création de post
│   ├── edit.html            # Formulaire de modification de post
│   ├── login.html           # Page de connexion
│   ├── register.html        # Page d'inscription
│   └── error.html           # Pages d'erreur (404, 403, 500)
│
├── static/
│   └── css/
│       ├── base.css         # Reset, variables, navbar, boutons (toutes les pages)
│       ├── index.css        # Page d'accueil : posts, filtres, sidebar, votes
│       ├── post.css         # Page de détail : commentaires, actions
│       ├── auth.css         # Pages login et register
│       ├── form.css         # Pages create et edit
│       └── error.css        # Page d'erreur

```

---

## Installation

### Prérequis

- [Go 1.21+](https://go.dev/dl/)
- GCC (requis par `go-sqlite3`) :
  ```bash
  # Ubuntu / Debian
  sudo apt install build-essential

  # macOS
  xcode-select --install
  ```

### Cloner le projet

```bash
git clone https://github.com/Lodgia/Forum_Greyddit.git
cd Forum_Greyddit
```

### Installer les dépendances

```bash
go mod tidy
```

### Lancer le serveur

```bash
go run main.go
```

Le serveur démarre sur **http://localhost:5050**.

La base de données `forum.db` est créée automatiquement au premier lancement avec les tables et les 5 catégories par défaut.

---

## Utilisation


### Accéder au site

| URL | Description |
|-----|-------------|
| `http://localhost:5050/` | Page d'accueil |
| `http://localhost:5050/register` | Inscription |
| `http://localhost:5050/login` | Connexion |
| `http://localhost:5050/post/create` | Créer un post (connecté) |
| `http://localhost:5050/category/{slug}` | Posts par catégorie |
| `http://localhost:5050/post/{id}` | Détail d'un post |

---

## Base de données

Le fichier `forum.db` contient 8 tables :

| Table | Contenu |
|-------|---------|
| `users` | Comptes utilisateurs |
| `categories` | Catégories des posts |
| `posts` | Les sujets du forum |
| `post_categories` | Liaison posts ↔ catégories |
| `comments` | Commentaires sur les posts |
| `post_likes` | Votes sur les posts (+1 / -1) |
| `comment_likes` | Votes sur les commentaires (+1 / -1) |
| `sessions` | Sessions de connexion actives |

Les catégories créées par défaut sont : **Vie étudiante**, **Projets**, **Bons plans**, **Détente**, **Entraide**.

---

## Sécurité

| Mesure | Détail |
|--------|--------|
| **Bcrypt** | Les mots de passe sont hashés avant stockage, jamais en clair |
| **Cookie HttpOnly** | Le cookie de session est inaccessible en JavaScript |
| **Sessions en base** | Chaque session a une expiration de 24h, nettoyées automatiquement |
| **UUID** | Les identifiants de sessions et d'entités sont des UUID aléatoires |
| **RequireAuth** | Les routes sensibles redirigent vers `/login` si non connecté |
| **Vérification propriétaire** | Impossible de modifier ou supprimer les ressources d'un autre utilisateur |
| **Paramètres SQL `?`** | Protection contre les injections SQL |
| **Clés étrangères** | `ON DELETE CASCADE` assure l'intégrité des données |

---

## Technologies utilisées

| Technologie | Usage |
|-------------|-------|
| [Go](https://go.dev/) | Langage principal, serveur HTTP |
| [SQLite](https://www.sqlite.org/) | Base de données embarquée |
| [go-sqlite3](https://github.com/mattn/go-sqlite3) | Driver SQLite pour Go |
| [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) | Hashage des mots de passe |
| [google/uuid](https://github.com/google/uuid) | Génération d'identifiants uniques |
| HTML / CSS | Templates et styles |

Aucun framework Go externe — uniquement la librairie standard.

---

## Auteurs

Projet réalisé dans le cadre du cours **Ynov Informatique — PROJET INFRA**.

- **EWEN**
- **TIMOTHE**
- **VAN**