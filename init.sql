CREATE TABLE `auth_service`(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL
);

CREATE TABLE `car`(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `license_plate` VARCHAR(255) NOT NULL DEFAULT '',
    `user_id` BIGINT UNSIGNED NOT NULL,
    `model_id` BIGINT UNSIGNED NOT NULL,
    `year` INT NOT NULL
);

CREATE TABLE `car_make`(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE `car_category`(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL UNIQUE,
    `passenger_count` INT NOT NULL
);

CREATE TABLE `car_model`(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `category_id` BIGINT UNSIGNED NOT NULL,
    `make_id` BIGINT UNSIGNED NOT NULL,
    `name` VARCHAR(255) NOT NULL
);

CREATE TABLE `ride_passenger`(
    `ride_id` BIGINT UNSIGNED NOT NULL,
    `passenger_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY(`ride_id`, `passenger_id`)

);

CREATE TABLE `auth`(
    `user_id` BIGINT UNSIGNED NOT NULL,
    `service_id` BIGINT UNSIGNED NOT NULL,
    `token` VARCHAR(255) NOT NULL,
    PRIMARY KEY(`user_id`, `service_id`)
);

CREATE TABLE `ride`(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `owner_user_id` BIGINT UNSIGNED NOT NULL,
    `vehicle_id` BIGINT UNSIGNED NOT NULL,
    `start_date` DATETIME NOT NULL,
    `start_city` VARCHAR(255) NOT NULL,
    `start_address` VARCHAR(255) NOT NULL,
    `end_city` VARCHAR(255) NOT NULL,
    `end_address` VARCHAR(255) NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP()
);

CREATE TABLE `chat_message`(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `ride_id` BIGINT UNSIGNED NOT NULL,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `message` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY(`id`, `ride_id`)
);

CREATE TABLE `user`(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `settings` JSON NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP()
);

CREATE TABLE `user_feedback`(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `owner_user_id` BIGINT UNSIGNED NOT NULL,
    `ride_id` BIGINT UNSIGNED NOT NULL,
    `score` INT NOT NULL,
    `message` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP()
);

-- Foreign Key Constraints
ALTER TABLE `ride` ADD CONSTRAINT `ride_owner_user_id_foreign` FOREIGN KEY(`owner_user_id`) REFERENCES `user`(`id`);
ALTER TABLE `user_feedback` ADD CONSTRAINT `user_feedback_ride_id_foreign` FOREIGN KEY(`ride_id`) REFERENCES `ride`(`id`);
ALTER TABLE `ride` ADD CONSTRAINT `ride_vehicle_id_foreign` FOREIGN KEY(`vehicle_id`) REFERENCES `car`(`id`);
ALTER TABLE `ride_passenger` ADD CONSTRAINT `ride_passenger_passenger_id_foreign` FOREIGN KEY(`passenger_id`) REFERENCES `user`(`id`);
ALTER TABLE `car_model` ADD CONSTRAINT `car_model_make_id_foreign` FOREIGN KEY(`make_id`) REFERENCES `car_make`(`id`);
ALTER TABLE `auth` ADD CONSTRAINT `auth_user_id_foreign` FOREIGN KEY(`user_id`) REFERENCES `user`(`id`);
ALTER TABLE `ride_passenger` ADD CONSTRAINT `ride_passenger_ride_id_foreign` FOREIGN KEY(`ride_id`) REFERENCES `ride`(`id`);
ALTER TABLE `car` ADD CONSTRAINT `car_user_id_foreign` FOREIGN KEY(`user_id`) REFERENCES `user`(`id`);
ALTER TABLE `car` ADD CONSTRAINT `car_model_id_foreign` FOREIGN KEY(`model_id`) REFERENCES `car_model`(`id`);
ALTER TABLE `auth` ADD CONSTRAINT `auth_service_foreign` FOREIGN KEY(`service_id`) REFERENCES `auth_service`(`id`);
ALTER TABLE `chat_message` ADD CONSTRAINT `chat_message_user_id_foreign` FOREIGN KEY(`user_id`) REFERENCES `user`(`id`);
ALTER TABLE `car_model` ADD CONSTRAINT `car_model_category_id_foreign` FOREIGN KEY(`category_id`) REFERENCES `car_category`(`id`);
ALTER TABLE `user_feedback` ADD CONSTRAINT `user_feedback_owner_user_id_foreign` FOREIGN KEY(`owner_user_id`) REFERENCES `user`(`id`);

INSERT INTO `user` (`email`, `name`, `password`, `settings`) VALUES
('tom@gmail.com', 'Tom Tommy', 'password1', '{}'),
('jerry@gmail.com', 'Jerry Jefferson', 'password2', '{}'),
('spike@gmail.com', 'Spike Spiky','password3', '{}'),
('tyke@gmail.com', 'Tyke Tyson','password4', '{}'),
('butch@gmail.com', 'Butch Butcherson', 'password5', '{}'),
('lightning@gmail.com', 'Lightning McQueen', 'password6', '{}'),
('tuffy@gmail.com', 'Tuffy Tufferson', 'password7', '{}'),
('muscles@gmail.com', 'Muscles Muscly', 'password8', '{}'),
('quacker@gmail.com', 'Tuffy Tufferson', 'password9', '{}'),
('nibbles@gmail.com', 'Nibbles Nibbleson', 'password10', '{}'),
('toodles@gmail.com', 'Toodles Toodleson', 'password11', '{}'),
('mammy@gmail.com', 'Mammy Mommy', 'password12', '{}'),
('george@gmail.com', 'George Washington', 'password13', '{}'),
('joan@gmail.com', 'Joan Joahnson', 'password14', '{}'),
('jeannie@gmail.com', 'Jeannie Marrie', 'password15', '{}'),
('goldie@gmail.com', 'Goldie Gulderson', 'password16', '{}'),
('fluff@gmail.com', 'Fluff Flufferson', 'password17', '{}'),
('meathead@gmail.com', 'Meathead Metalicca', 'password18', '{}'),
('cuckoo@gmail.com', 'Cuckoo Cucumber', 'password19', '{}'),
('puddy@gmail.com', 'Puddy Pudgy', 'password20', '{}');

INSERT INTO `car_make` (`name`) VALUES
('Toyota'),
('Honda'),
('Ford'),
('Chevrolet'),
('Nissan'),
('BMW'),
('Mercedes'),
('Volkswagen'),
('Audi'),
('Hyundai'),
('Kia'),
('Mazda'),
('Subaru'),
('Lexus'),
('Jaguar'),
('Porsche'),
('Ferrari'),
('Lamborghini'),
('Bentley'),
('Rolls-Royce');

INSERT INTO `car_category` (`name`, `passenger_count`) VALUES
('Sedan', 5),
('SUV', 7),
('Truck', 3),
('Coupe', 4),
('Convertible', 4),
('Minivan', 7),
('Hatchback', 5),
('Wagon', 5),
('Sports Car', 2),
('Diesel', 5),
('Electric', 5),
('Hybrid', 5),
('Luxury', 5),
('Off-Road', 5),
('Pickup', 5),
('Van', 5),
('Compact', 5),
('Subcompact', 5),
('Crossover', 5),
('Roadster', 2);

INSERT INTO `car_model` (`category_id`, `make_id`, `name`) VALUES
(1, 1, 'Camry'),
(1, 1, 'Corolla'),
(1, 1, 'Prius'),
(1, 1, 'Avalon'),
(1, 1, 'Yaris'),
(9, 1, 'Supra'),
(2, 1, 'Highlander'),
(2, 1, 'RAV4'),
(3, 1, 'Tacoma'),
(3, 1, 'Tundra'),
(1, 2, 'Civic'),
(1, 2, 'Accord'),
(1, 2, 'Fit'),
(1, 2, 'Insight'),
(2, 2, 'CR-V'),
(2, 2, 'HR-V'),
(2, 2, 'Passport'),
(2, 2, 'Pilot'),
(3, 2, 'Ridgeline'),
(6, 2, 'Odyssey'),
(3, 3, 'F-150'),
(9, 3, 'Mustang'),
(2, 3, 'Explorer'),
(2, 3, 'Escape'),
(2, 3, 'Edge'),
(2, 3, 'Expedition'),
(3, 3, 'Ranger'),
(2, 3, 'Bronco'),
(1, 3, 'Fusion'),
(1, 3, 'Taurus'),
(9, 4, 'Camaro'),
(1, 4, 'Malibu'),
(1, 4, 'Impala'),
(1, 4, 'Cruze'),
(1, 4, 'Spark'),
(1, 4, 'Sonic'),
(2, 4, 'Trax'),
(2, 4, 'Equinox'),
(2, 4, 'Blazer'),
(2, 4, 'Traverse'),
(1, 5, 'Altima'),
(1, 5, 'Sentra'),
(1, 5, 'Maxima'),
(1, 5, 'Versa'),
(11, 5, 'Leaf'),
(2, 5, 'Juke'),
(2, 5, 'Rogue'),
(2, 5, 'Murano'),
(2, 5, 'Pathfinder'),
(2, 5, 'Armada'),
(1, 6, '3 Series'),
(1, 6, '5 Series'),
(1, 6, '7 Series'),
(2, 6, 'X1'),
(2, 6, 'X3'),
(2, 6, 'X5'),
(2, 6, 'X7'),
(9, 6, 'Z4'),
(11, 6, 'i3'),
(11, 6, 'i8'),
(1, 7, 'C-Class'),
(1, 7, 'E-Class'),
(1, 7, 'S-Class'),
(2, 7, 'GLA'),
(2, 7, 'GLC'),
(2, 7, 'GLE'),
(2, 7, 'GLS'),
(1, 7, 'A-Class'),
(1, 7, 'B-Class'),
(1, 7, 'CLA'),
(1, 8, 'Golf'),
(1, 8, 'Jetta'),
(1, 8, 'Passat'),
(9, 8, 'Beetle'),
(2, 8, 'Tiguan'),
(2, 8, 'Atlas'),
(2, 8, 'Touareg'),
(1, 8, 'Arteon'),
(1, 8, 'CC'),
(1, 8, 'Eos'),
(1, 9, 'A4'),
(1, 9, 'A6'),
(1, 9, 'A8'),
(2, 9, 'Q3'),
(2, 9, 'Q5'),
(2, 9, 'Q7'),
(2, 9, 'Q8'),
(9, 9, 'TT'),
(9, 9, 'R8'),
(1, 9, 'S4'),
(1, 10, 'Elantra'),
(1, 10, 'Sonata'),
(1, 10, 'Accent'),
(9, 10, 'Veloster'),
(2, 10, 'Kona'),
(2, 10, 'Tucson'),
(2, 10, 'Santa Fe'),
(2, 10, 'Palisade'),
(11, 10, 'Ioniq'),
(11, 10, 'Nexo'),
(1, 11, 'Soul'),
(1, 11, 'Forte'),
(1, 11, 'Optima'),
(9, 11, 'Stinger'),
(1, 11, 'Rio'),
(2, 11, 'Seltos'),
(2, 11, 'Sportage'),
(2, 11, 'Sorento'),
(2, 11, 'Telluride'),
(11, 11, 'Niro'),
(2, 12, 'CX-5'),
(2, 12, 'CX-3'),
(2, 12, 'CX-9'),
(1, 12, 'Mazda3'),
(1, 12, 'Mazda6'),
(9, 12, 'MX-5 Miata'),
(2, 12, 'CX-30'),
(2, 12, 'CX-50'),
(2, 12, 'CX-60'),
(2, 12, 'CX-90'),
(2, 13, 'Outback'),
(2, 13, 'Forester'),
(1, 13, 'Impreza'),
(1, 13, 'Legacy'),
(2, 13, 'Crosstrek'),
(2, 13, 'Ascent'),
(9, 13, 'BRZ'),
(9, 13, 'WRX'),
(2, 13, 'XV'),
(2, 13, 'Levorg'),
(2, 14, 'RX'),
(2, 14, 'NX'),
(2, 14, 'UX'),
(2, 14, 'GX'),
(2, 14, 'LX'),
(1, 14, 'ES'),
(1, 14, 'IS'),
(1, 14, 'GS'),
(1, 14, 'LS'),
(9, 14, 'LC'),
(9, 15, 'F-Type'),
(1, 15, 'XE'),
(1, 15, 'XF'),
(1, 15, 'XJ'),
(2, 15, 'E-Pace'),
(2, 15, 'F-Pace'),
(11, 15, 'I-Pace'),
(9, 15, 'XK'),
(9, 15, 'XKR'),
(1, 15, 'XFR'),
(9, 16, '911'),
(2, 16, 'Cayenne'),
(2, 16, 'Macan'),
(1, 16, 'Panamera'),
(11, 16, 'Taycan'),
(9, 16, 'Boxster'),
(9, 16, 'Cayman'),
(9, 16, 'Carrera'),
(9, 16, 'Turbo'),
(9, 16, 'Spyder'),
(9, 17, '488'),
(9, 17, '812 Superfast'),
(9, 17, 'Portofino'),
(9, 17, 'Roma'),
(11, 17, 'SF90 Stradale'),
(9, 17, 'F8 Tributo'),
(2, 17, 'GTC4Lusso'),
(9, 17, 'LaFerrari'),
(9, 17, 'Monza SP1'),
(9, 17, 'Monza SP2'),
(9, 18, 'Huracan'),
(9, 18, 'Aventador'),
(2, 18, 'Urus'),
(9, 18, 'Gallardo'),
(9, 18, 'Murcielago'),
(9, 18, 'Diablo'),
(9, 18, 'Countach'),
(9, 18, 'Reventon'),
(9, 18, 'Veneno'),
(9, 18, 'Sesto Elemento'),
(9, 19, 'Continental'),
(1, 19, 'Flying Spur'),
(2, 19, 'Bentayga'),
(1, 19, 'Mulsanne'),
(1, 19, 'Arnage'),
(1, 19, 'Azure'),
(1, 19, 'Brooklands'),
(1, 19, 'Turbo R'),
(1, 19, 'Eight'),
(1, 19, 'T Series'),
(9, 20, 'Phantom'),
(1, 20, 'Ghost'),
(9, 20, 'Wraith'),
(9, 20, 'Dawn'),
(2, 20, 'Cullinan'),
(1, 20, 'Silver Shadow'),
(1, 20, 'Silver Spirit'),
(1, 20, 'Silver Spur'),
(1, 20, 'Corniche'),
(1, 20, 'Camargue');

INSERT INTO `car` (`user_id`, `year`, `model_id`) VALUES
(1, 2020, 1),
(2, 2019, 2),
(3, 2018, 3),
(3, 2021, 4);

-- Insert rides
INSERT INTO `ride` (`owner_user_id`, `vehicle_id`, `start_date`, `start_city`, `start_address`, `end_city`, `end_address`) VALUES
(1, 1, '2023-10-01 08:00:00', 'New York', '123 Main St', 'Boston', '456 Elm St'),
(1, 1, '2023-10-02 09:00:00', 'Los Angeles', '789 Oak St', 'San Francisco', '101 Pine St'),
(2, 2, '2023-10-03 10:00:00', 'Chicago', '202 Maple St', 'Detroit', '303 Birch St'),
(2, 3, '2023-10-04 11:00:00', 'Houston', '404 Cedar St', 'Dallas', '505 Walnut St'),
(3, 4, '2023-10-05 12:00:00', 'Phoenix', '606 Spruce St', 'Tucson', '707 Fir St'),
(3, 4, '2023-10-06 13:00:00', 'Philadelphia', '808 Ash St', 'Pittsburgh', '909 Poplar St'),
(3, 4, '2023-10-07 14:00:00', 'San Antonio', '1010 Willow St', 'Austin', '1111 Cypress St'),
(3, 4, '2023-10-08 15:00:00', 'San Diego', '1212 Redwood St', 'Las Vegas', '1313 Palm St'),
(3, 4, '2023-10-09 16:00:00', 'Dallas', '1414 Magnolia St', 'Houston', '1515 Dogwood St'),
(3, 4, '2023-10-10 17:00:00', 'San Jose', '1616 Cherry St', 'Sacramento', '1717 Peach St');

-- Insert ride passengers
INSERT INTO `ride_passenger` (`ride_id`, `passenger_id`) VALUES
(1, 2), (1, 3),
(2, 4), (2, 5),
(3, 6), (3, 7);

-- Insert user feedback
INSERT INTO `user_feedback` (`owner_user_id`, `ride_id`, `score`, `message`) VALUES
(2, 1, 5, 'Great ride!'),
(3, 1, 4, 'Very comfortable.'),
(4, 2, 5, 'Excellent driver.'),
(5, 2, 3, 'Good, but could be better.'),
(6, 3, 5, 'Fantastic experience.'),
(7, 3, 4, 'Nice and smooth.'),
(8, 4, 5, 'Loved it!'),
(9, 4, 4, 'Pretty good.'),
(10, 5, 5, 'Amazing ride.'),
(1, 5, 3, 'It was okay.');

-- Insert chat messages
INSERT INTO `chat_message` (`ride_id`, `user_id`, `message`, `created_at`) VALUES
(1, 1, 'Ready to go?', '2023-10-01 07:50:00'),
(1, 2, 'Yes, I am.', '2023-10-01 07:51:00'),
(2, 2, 'On my way.', '2023-10-02 08:50:00'),
(2, 4, 'Great, see you soon.', '2023-10-02 08:51:00'),
(3, 3, 'Leaving now.', '2023-10-03 09:50:00'),
(3, 6, 'Got it.', '2023-10-03 09:51:00'),
(4, 4, 'Starting the ride.', '2023-10-04 10:50:00'),
(4, 8, 'Okay.', '2023-10-04 10:51:00'),
(5, 5, 'Heading out.', '2023-10-05 11:50:00'),
(5, 10, 'See you soon.', '2023-10-05 11:51:00'),
(6, 6, 'On the way.', '2023-10-06 12:50:00'),
(6, 2, 'Alright.', '2023-10-06 12:51:00'),
(7, 7, 'Leaving now.', '2023-10-07 13:50:00'),
(7, 4, 'Got it.', '2023-10-07 13:51:00'),
(8, 8, 'Starting the ride.', '2023-10-08 14:50:00'),
(8, 6, 'Okay.', '2023-10-08 14:51:00'),
(9, 9, 'Heading out.', '2023-10-09 15:50:00'),
(9, 8, 'See you soon.', '2023-10-09 15:51:00'),
(10, 10, 'On the way.', '2023-10-10 16:50:00'),
(10, 1, 'Alright.', '2023-10-10 16:51:00');