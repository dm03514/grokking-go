#!/usr/bin/env sh

psql -c "create database deposits"

psql -c "create schema deposits"

psql deposits -c "CREATE TABLE deposits (
  id SERIAL,
  transaction_id INT,
  account_number INT NOT NULL,
  amount_cents INT NOT NULL ,
  PRIMARY KEY (transaction_id)
)"