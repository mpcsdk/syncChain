--
-- PostgreSQL database dump
--

-- Dumped from database version 12.3 (Debian 12.3-1.pgdg100+1)
-- Dumped by pg_dump version 12.16 (Ubuntu 12.16-0ubuntu0.20.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: chain_data; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chain_data (
    chain_id bigint NOT NULL,
    height bigint NOT NULL,
    block_hash character varying(255) NOT NULL,
    ts bigint NOT NULL,
    tx_hash character varying(255) NOT NULL,
    tx_idx integer NOT NULL,
    log_idx integer NOT NULL,
    "from" character varying(255) NOT NULL,
    "to" character varying(255) NOT NULL,
    contract character varying(255) NOT NULL,
    value character varying(255) NOT NULL,
    gas character varying(255) NOT NULL,
    gas_price character varying(255) NOT NULL,
    nonce bigint
);


ALTER TABLE public.chain_data OWNER TO postgres;

--
-- Name: fromtscontractid; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fromtscontractid ON public.chain_data USING btree (ts DESC NULLS LAST, "from", contract, chain_id);


--
-- Name: hashtxidxlogidx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX hashtxidxlogidx ON public.chain_data USING btree (chain_id, tx_hash, tx_idx, log_idx);


--
-- Name: totscontractid; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX totscontractid ON public.chain_data USING btree (ts DESC NULLS LAST, "to", contract, chain_id);


--
-- Name: tscontractid; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX tscontractid ON public.chain_data USING btree (chain_id, ts DESC NULLS LAST, contract);


--
-- PostgreSQL database dump complete
--

