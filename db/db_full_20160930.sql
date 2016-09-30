--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.0
-- Dumped by pg_dump version 9.6.0

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: twitter_users; Type: TABLE; Schema: public; Owner: perelin
--

CREATE TABLE twitter_users (
    user_id text,
    screen_name text,
    oauth_token text,
    oauth_token_secret text,
    create_date timestamp without time zone
);


ALTER TABLE twitter_users OWNER TO perelin;

--
-- Data for Name: twitter_users; Type: TABLE DATA; Schema: public; Owner: perelin
--

COPY twitter_users (user_id, screen_name, oauth_token, oauth_token_secret, create_date) FROM stdin;
7559392	perelin	7559392-ZGAbEunH3GIZuYy0DNLCB62H5uAX4rxmHyQTljws7a	5e0VWLJomshYsKm7SGxZDQjnNJ9MeUeb7c96OyL3TrJir	2016-09-30 23:12:18.870417
\.


--
-- PostgreSQL database dump complete
--

