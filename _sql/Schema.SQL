-- Aller dans la DB :  & "C:\Users\natha\tools\scalingo.exe" --app projetyboost mysql-console
-- Execute fichier SQL : Get-Content .\Schema.sql | & "C:\Users\natha\tools\scalingo.exe" --app projetyboost mysql-console
-- Voir la base de données : SELECT DATABASE();
-- Tunnel :  & "C:\Users\natha\tools\scalingo.exe" --app projetyboost --region osc-fr1 db-tunnel SCALINGO_MYSQL_URL



-- Daily champ table
CREATE TABLE IF NOT EXISTS daily_champion (
  day DATE PRIMARY KEY,
  champion_id VARCHAR(50) NOT NULL
);

-- Users auth table
CREATE TABLE IF NOT EXISTS users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  pseudo VARCHAR(100) NOT NULL UNIQUE,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP

);



exit 