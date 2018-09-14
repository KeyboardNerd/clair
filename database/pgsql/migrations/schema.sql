
CREATE TABLE IF NOT EXISTS namespace (
id SERIAL PRIMARY KEY,
name TEXT NULL,
version_format TEXT,
UNIQUE (name, version_format));
CREATE INDEX ON namespace(name);


CREATE TABLE IF NOT EXISTS feature (
id SERIAL PRIMARY KEY,
name TEXT NOT NULL,
version TEXT NOT NULL,
version_format TEXT NOT NULL,
UNIQUE (name, version, version_format));
CREATE INDEX ON feature(name);

CREATE TABLE IF NOT EXISTS namespaced_feature (
id SERIAL PRIMARY KEY,
namespace_id INT REFERENCES namespace,
feature_id INT REFERENCES feature,
UNIQUE (namespace_id, feature_id));


CREATE TYPE detector_type AS ENUM ('namespace', 'feature');


CREATE TABLE IF NOT EXISTS detector (
id SERIAL PRIMARY KEY,
name TEXT NOT NULL,
version TEXT NOT NULL,
type detector_type NOT NULL,
UNIQUE (name, version, type));


CREATE TABLE IF NOT EXISTS layer(
id SERIAL PRIMARY KEY,
hash TEXT NOT NULL UNIQUE);

CREATE TABLE IF NOT EXISTS layer_detector(
id SERIAL PRIMARY KEY,
layer_id INT REFERENCES layer ON DELETE CASCADE,
detector_id INT REFERENCES detector ON DELETE CASCADE,
UNIQUE(layer_id, detector_id));
CREATE INDEX ON layer_detector(layer_id);

CREATE TABLE IF NOT EXISTS layer_feature (
id SERIAL PRIMARY KEY,
layer_id INT REFERENCES layer ON DELETE CASCADE, 
feature_id INT REFERENCES feature ON DELETE CASCADE,
detector_id INT REFERENCES detector ON DELETE CASCADE,
UNIQUE (layer_id, feature_id));
CREATE INDEX ON layer_feature(layer_id);

CREATE TABLE IF NOT EXISTS layer_namespace (
id SERIAL PRIMARY KEY,
layer_id INT REFERENCES layer ON DELETE CASCADE,
namespace_id INT REFERENCES namespace ON DELETE CASCADE,
detector_id INT REFERENCES detector ON DELETE CASCADE,
UNIQUE (layer_id, namespace_id));
CREATE INDEX ON layer_namespace(layer_id);



CREATE TABLE IF NOT EXISTS ancestry (
id SERIAL PRIMARY KEY,
name TEXT NOT NULL UNIQUE);

CREATE TABLE IF NOT EXISTS ancestry_layer (
id SERIAL PRIMARY KEY,
ancestry_id INT REFERENCES ancestry ON DELETE CASCADE,
ancestry_index INT NOT NULL,
layer_id INT REFERENCES layer ON DELETE RESTRICT,
UNIQUE (ancestry_id, ancestry_index));
CREATE INDEX ON ancestry_layer(ancestry_id);

CREATE TABLE IF NOT EXISTS ancestry_feature(
id SERIAL PRIMARY KEY,
ancestry_layer_id INT REFERENCES ancestry_layer ON DELETE CASCADE,
namespaced_feature_id INT REFERENCES namespaced_feature ON DELETE CASCADE,
feature_detector_id INT REFERENCES detector ON DELETE CASCADE,
namespace_detector_id INT REFERENCES detector ON DELETE CASCADE,
UNIQUE (ancestry_layer_id, namespaced_feature_id));

CREATE TABLE IF NOT EXISTS ancestry_detector(
id SERIAL PRIMARY KEY,
ancestry_id INT REFERENCES layer ON DELETE CASCADE,
detector_id INT REFERENCES detector ON DELETE CASCADE,
UNIQUE(ancestry_id, detector_id));
CREATE INDEX ON ancestry_detector(ancestry_id);


CREATE TYPE severity AS ENUM ('Unknown', 'Negligible', 'Low', 'Medium', 'High', 'Critical', 'Defcon1');

CREATE TABLE IF NOT EXISTS vulnerability (
id SERIAL PRIMARY KEY,
namespace_id INT NOT NULL REFERENCES Namespace,
name TEXT NOT NULL,
description TEXT NULL,
link TEXT NULL,
severity severity NOT NULL,
metadata TEXT NULL,
created_at TIMESTAMP WITH TIME ZONE,
deleted_at TIMESTAMP WITH TIME ZONE NULL);
CREATE INDEX ON vulnerability(namespace_id, name);
CREATE INDEX ON vulnerability(namespace_id);

CREATE TABLE IF NOT EXISTS vulnerability_affected_feature (
id SERIAL PRIMARY KEY, 
vulnerability_id INT NOT NULL REFERENCES vulnerability ON DELETE CASCADE,
feature_name TEXT NOT NULL,
affected_version TEXT,
fixedin TEXT);
CREATE INDEX ON vulnerability_affected_feature(vulnerability_id, feature_name);

CREATE TABLE IF NOT EXISTS vulnerability_affected_namespaced_feature(
id SERIAL PRIMARY KEY,
vulnerability_id INT NOT NULL REFERENCES vulnerability ON DELETE CASCADE,
namespaced_feature_id INT NOT NULL REFERENCES namespaced_feature ON DELETE CASCADE,
added_by INT NOT NULL REFERENCES vulnerability_affected_feature ON DELETE CASCADE,
UNIQUE (vulnerability_id, namespaced_feature_id));
CREATE INDEX ON vulnerability_affected_namespaced_feature(namespaced_feature_id);


CREATE TABLE IF NOT EXISTS KeyValue (
id SERIAL PRIMARY KEY,
key TEXT NOT NULL UNIQUE,
VALUE TEXT);

CREATE TABLE IF NOT EXISTS Lock (
id SERIAL PRIMARY KEY,
name VARCHAR(64) NOT NULL UNIQUE,
owner VARCHAR(64) NOT NULL,
until TIMESTAMP WITH TIME ZONE);
CREATE INDEX ON Lock (owner);


CREATE TABLE IF NOT EXISTS Vulnerability_Notification (
id SERIAL PRIMARY KEY,
name VARCHAR(64) NOT NULL UNIQUE,
created_at TIMESTAMP WITH TIME ZONE,
notified_at TIMESTAMP WITH TIME ZONE NULL,
deleted_at TIMESTAMP WITH TIME ZONE NULL,
old_vulnerability_id INT NULL REFERENCES Vulnerability ON DELETE CASCADE,
new_vulnerability_id INT NULL REFERENCES Vulnerability ON DELETE CASCADE);
CREATE INDEX ON Vulnerability_Notification (notified_at);