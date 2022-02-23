SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";

--
-- Database: `flow-user`
--

CREATE DATABASE IF NOT EXISTS `flow-user` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE `flow-user`;

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL UNIQUE,
  `password` varchar(255) NOT NULL,
  PRIMARY KEY (id)
);

--
-- Table structure for table `github_oauth2_tokens`
--

CREATE TABLE `github_oauth2_tokens` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint UNSIGNED NOT NULL UNIQUE,
  `access_token` varchar(255) NOT NULL,
  `owner_id` varchar(255) NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

--
-- Table structure for table `google_oauth2_tokens`
--

CREATE TABLE `google_oauth2_tokens` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint UNSIGNED NOT NULL UNIQUE,
  `access_token` varchar(255) NOT NULL,
  `access_token_expire_in` datetime NOT NULL,
  `refresh_token` varchar(255) NOT NULL,
  `owner_id` varchar(255) NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

--
-- Table structure for table `twitter_oauth2_tokens`
--

CREATE TABLE `twitter_oauth2_tokens` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint UNSIGNED NOT NULL UNIQUE,
  `access_token` varchar(255) NOT NULL,
  `access_token_expire_in` datetime NOT NULL,
  `refresh_token` varchar(255) NOT NULL,
  `refresh_token_expire_in` datetime NOT NULL,
  `owner_id` varchar(255) NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);