# UHC LoLdle

## Description

UHC LoLdle est une application web inspirée de l’univers de League of Legends.
Le projet propose plusieurs défis quotidiens autour des champions, des objets et des sorts du jeu.

L’objectif est de permettre à un utilisateur connecté de tester ses connaissances grâce à différents challenges avec indices progressifs, historique des essais, système de streak et interface personnalisée.

## Fonctionnalités principales

### Authentification
- Inscription d’un nouvel utilisateur
- Connexion sécurisée avec mot de passe hashé
- Déconnexion
- Gestion de session par cookie

### Profil utilisateur
- Affichage du pseudo de l'utilisateur
- Affichage de la date d’inscription
- Affichage de la série actuelle
- Affichage de la meilleure série
- Affichage du nombre total de jours joués

### Challenge Champion
- Un champion du jour est sélectionné et enregistré en base
- Le joueur doit retrouver le champion à partir de son lore
- Le nom du champion est masqué dans le texte
- Les indices se débloquent au fil des essais
- Le lore complet est révélé une fois la bonne réponse trouvée

### Challenge Item
- Un item du jour est sélectionné et enregistré en base
- Le joueur doit deviner l’objet à partir de ses statistiques
- Le prix est révélé après plusieurs essais
- Les composants de l’objet sont affichés après davantage d’essais

### Challenge Spell
- Un sort quotidien est sélectionné à partir d’un champion
- Le joueur doit retrouver à quel champion appartient le sort
- Des indices comme le nom du sort et son image se débloquent progressivement
- Une seconde étape demande d’identifier le slot du sort : Q, W, E ou R

### Fonctionnalités supplémentaires
- Sauvegarde locale de l’état de progression avec localStorage
- Compte à rebours jusqu’au prochain défi quotidien
- Désactivation ou réactivation des indices
- Modales de règles pour chaque challenge
- Interface responsive avec sidebar fixe

## Technologies utilisées

### Backend
- Go
- net/http
- html/template
- database/sql
- MySQL
- bcrypt pour le hash des mots de passe

### Frontend
- HTML
- CSS
- JavaScript

### API externe
- Riot Data Dragon API

### Déploiement
- GitHub
- Scalingo
- MySQL addon sur Scalingo

## Structure du projet

```
YboostAppProject/
├── cmd/
│   ├── main.go
│   ├── Handler.go
│   ├── Cookie.go
│   ├── RegisterSteps.go
│   ├── getInfo.go
│   ├── ChampionChallenge.go
│   ├── ItemChallenge.go
│   ├── SpellChallenge.go
│   ├── champStruct.go
│   ├── ItemStruct.go
│   └── SpellStruct.go
├── static/
│   ├── css/
│   ├── images/
│   └── js/
├── templates/
│   ├── HomePage.html
│   ├── LoginPage.html
│   ├── RegisterPage.html
│   ├── ProfilPage.html
│   ├── ChampPage.html
│   ├── ItemsPage.html
│   └── SpellPage.html
├── _sql/
│   └── Schema.SQL
├── go.mod
├── go.sum
└── Procfile
```

## Rôle des principaux fichiers

### main.go
Point d’entrée de l’application.
Initialise la connexion à la base de données, charge les données utiles au démarrage et enregistre les routes HTTP.

### Handler.go
Contient les handlers principaux :
- page d’accueil
- inscription
- connexion
- profil
- pages des challenges
- vérification des réponses
- déconnexion

### Cookie.go
Gère les sessions utilisateur, la récupération du compte courant et la logique de mise à jour du streak.

### RegisterSteps.go
Contient la validation des champs d’inscription et la logique de hash du mot de passe.

### getInfo.go
Récupère les données depuis l’API Riot :
- version de Data Dragon
- champions
- lore complet
- items
- composants d’items

### ChampionChallenge.go
Gère la logique du challenge champion :
- sélection quotidienne
- persistance en base
- comparaison des réponses

### ItemChallenge.go
Gère la logique du challenge item :
- sélection quotidienne
- persistance en base
- comparaison des réponses

### SpellChallenge.go
Gère la logique du challenge sort :
- sélection quotidienne d’un champion et d’un slot
- récupération du sort concerné
- préparation des données pour le challenge

### Fichiers de structures
- champStruct.go
- ItemStruct.go
- SpellStruct.go

Ces fichiers regroupent les structures utilisées pour parser les réponses JSON de l’API Riot et manipuler les données métier du projet.

### common.js
Contient les fonctions JavaScript réutilisables entre plusieurs pages :
- filtrage des listes
- rendu des suggestions
- rendu des statistiques
- composants d’items
- création des boutons de slot
- gestion des modales
- compte à rebours quotidien

## Base de données

Le projet s’appuie sur une base MySQL.

### Tables principales
- users
- daily_champion
- daily_item
- daily_spell

### Table users
Stocke les informations utilisateur :
- id
- email
- pseudo
- password_hash
- created_at
- streak
- best_streak
- day_played
- last_played

### Table daily_champion
Stocke le champion du jour.

### Table daily_item
Stocke l’item du jour.

### Table daily_spell
Stocke le champion et le slot du sort du jour.

## Accéssibilité sur navigateur 

Accéder à l'url : https://projetyboost.osc-fr1.scalingo.io/login 

Créer un compte et tester votre connaissance du jeu !



## Installation locale

### Prérequis
- Go installé
- MySQL installé
- Git
- Une base de données MySQL créée

### Étapes
1. Cloner le projet
2. Créer la base de données
3. Exécuter le script SQL de création
4. Configurer l’URL de connexion à la base
5. Installer les dépendances Go
6. Lancer l’application

### Commandes utiles

git clone https://github.com/nathlme/YboostAppProject.git

cd YboostAppProject

go mod tidy

go run ./cmd

## Lancement local

L’application est lancée depuis le dossier cmd.

Exemple :

go run ./cmd

Ensuite, ouvrir le navigateur sur l’adresse locale configurée par le serveur.

## Déploiement sur Scalingo

Le projet peut être déployé sur Scalingo avec :
- un dépôt Git
- un addon MySQL
- un Procfile
- les variables d’environnement nécessaires

### Étapes générales
1. Créer l’application sur Scalingo
2. Ajouter un addon MySQL
3. Configurer la variable de connexion à la base
4. Ajouter le remote Git Scalingo
5. Push le projet sur Scalingo

### Exemple de commandes

git remote add scalingo git@ssh.osc-fr1.scalingo.com:nom-de-l-app.git
git push scalingo main

## Fonctionnement général de l’application

### 1. Connexion utilisateur
L’utilisateur s’inscrit ou se connecte.

### 2. Chargement des données
Les données principales sont récupérées ou préparées :
- liste des champions
- liste des items
- défi quotidien

### 3. Affichage des pages
Les handlers injectent les données nécessaires dans les templates HTML.

### 4. Interaction utilisateur
Le joueur saisit une réponse dans un champ de recherche avec suggestions dynamiques.

### 5. Vérification serveur
Une requête POST est envoyée au backend pour vérifier la réponse.

### 6. Mise à jour du streak
Si la réponse est correcte, le streak est mis à jour en base selon les règles du jour.

### 7. Sauvegarde locale
La progression est conservée dans le navigateur avec localStorage.

## Sécurité mise en place

- Hash des mots de passe avec bcrypt
- Requêtes SQL paramétrées
- Gestion de session via cookie HttpOnly
- Redirection des routes protégées vers la connexion si l’utilisateur n’est pas authentifié

## Axes d’amélioration possibles

- Ajouter une protection CSRF
- Ajouter des statistiques plus avancées au profil
- Ajouter une meilleure gestion des erreurs réseau liées à l’API Riot
- Ajouter des tests automatisés

## Objectifs pédagogiques du projet

Ce projet permet de mettre en pratique :
- le développement d’une application web en Go
- le routage HTTP
- l’utilisation des templates HTML
- la manipulation du DOM en JavaScript
- l’utilisation d’une API externe
- la gestion d’une base de données relationnelle
- la gestion d’authentification et de session
- le déploiement d’une application web

## Auteur

Projet développé dans le cadre d’un projet web autour de l’univers de League of Legends - LAMARCHE Nathan

## Résumé

UHC LoLdle est une application web de défis quotidiens inspirée de League of Legends.
Elle combine backend Go, frontend HTML/CSS/JavaScript, base de données MySQL et API Riot pour proposer une expérience interactive avec authentification, sauvegarde de progression et défis variés.