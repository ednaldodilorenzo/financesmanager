CREATE TABLE `users` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password` varchar(100) NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `userscol_UNIQUE` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;

INSERT INTO `users` (`name`, `email`, `password`, `created_at`) VALUES ('Master', 'master@test.com', '"$2a$10$Cl84JkHPTsM0lWdgDOQxduCU6YRC9qR7glgDVIfJU0W2995sXLAGS"', curdate());