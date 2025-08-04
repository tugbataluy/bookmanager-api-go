--
-- PostgreSQL database dump
--

-- Dumped from database version 14.18 (Ubuntu 14.18-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.18 (Ubuntu 14.18-0ubuntu0.22.04.1)

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
-- Name: books; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.books (
    id integer NOT NULL,
    title character varying(255) NOT NULL,
    author character varying(255) NOT NULL,
    published_date date NOT NULL,
    edition integer DEFAULT 1 NOT NULL,
    description text,
    genre character varying(100),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.books OWNER TO postgres;

--
-- Name: books_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.books_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.books_id_seq OWNER TO postgres;

--
-- Name: books_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.books_id_seq OWNED BY public.books.id;


--
-- Name: collection_books; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.collection_books (
    collection_id integer NOT NULL,
    book_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.collection_books OWNER TO postgres;

--
-- Name: collections; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.collections (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.collections OWNER TO postgres;

--
-- Name: collections_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.collections_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.collections_id_seq OWNER TO postgres;

--
-- Name: collections_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.collections_id_seq OWNED BY public.collections.id;


--
-- Name: books id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.books ALTER COLUMN id SET DEFAULT nextval('public.books_id_seq'::regclass);


--
-- Name: collections id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.collections ALTER COLUMN id SET DEFAULT nextval('public.collections_id_seq'::regclass);


--
-- Data for Name: books; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.books (id, title, author, published_date, edition, description, genre, created_at, updated_at) FROM stdin;
1	Dune	Frank Herbert	1965-08-01	1	The first book in the Dune series	Science Fiction	2025-08-04 12:40:41.255674+03	2025-08-04 12:40:41.255674+03
2	Dune Messiah	Frank Herbert	1969-01-01	1	The second book in the Dune series	Science Fiction	2025-08-04 12:59:16.57765+03	2025-08-04 12:59:16.57765+03
3	Children of Dune	Frank Herbert	1976-01-01	1	The third book in the Dune series	Science Fiction	2025-08-04 13:00:19.376355+03	2025-08-04 13:00:19.376355+03
4	God Emperor of Dune	Frank Herbert	1981-01-01	1	The fourth book in the Dune series	Science Fiction	2025-08-04 13:01:17.54567+03	2025-08-04 13:01:17.54567+03
6	Chapterhouse: Dune	Frank Herbert	1985-01-01	1	The sixth book in the Dune series	Science Fiction	2025-08-04 13:02:35.723162+03	2025-08-04 13:02:35.723162+03
5	Heretics of Dune (Updated)	Frank Herbert	1984-01-01	2	The fifth book in the Dune series (My fav)	Science Fiction	2025-08-04 13:02:08.528665+03	2025-08-04 13:28:27.988437+03
8	The Go Programming Language	Alan A. A. Donovan & Brian W. Kernighan	2015-10-26	1	Authoritative resource for learning Go	Programming	2025-08-04 15:12:48.269931+03	2025-08-04 15:12:48.269931+03
9	Clean Code: A Handbook of Agile Software Craftsmanship	Robert C. Martin	2008-08-01	1	Principles for writing maintainable code	Software Engineering	2025-08-04 15:16:54.943853+03	2025-08-04 15:16:54.943853+03
11	Introduction to Algorithms	Thomas H. Cormen	2009-07-31	1	Comprehensive guide to algorithms	Computer Science	2025-08-04 15:25:18.029089+03	2025-08-04 15:25:18.029089+03
12	The Go Programming Language	Alan A. A. Donovan & Brian W. Kernighan	2015-10-26	1	Authoritative resource for learning Go	Programming	2025-08-04 15:57:24.42583+03	2025-08-04 15:57:24.42583+03
10	Design Patterns: Elements of Reusable Object-Oriented Software	Erich Gamma, Richard Helm, Ralph Johnson	1994-10-31	1	Classic solutions to common design problems	Software Engineering	2025-08-04 15:22:55.984542+03	2025-08-04 16:04:09.849435+03
\.


--
-- Data for Name: collection_books; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.collection_books (collection_id, book_id, created_at) FROM stdin;
1	1	2025-08-04 14:00:14.83678+03
1	3	2025-08-04 16:17:42.504292+03
\.


--
-- Data for Name: collections; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.collections (id, name, description, created_at, updated_at) FROM stdin;
1	Science Fiction Novels	A collection of sci-fi books.	2025-08-04 13:34:35.385499+03	2025-08-04 13:54:16.206639+03
\.


--
-- Name: books_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.books_id_seq', 13, true);


--
-- Name: collections_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.collections_id_seq', 3, true);


--
-- Name: books books_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.books
    ADD CONSTRAINT books_pkey PRIMARY KEY (id);


--
-- Name: collection_books collection_books_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.collection_books
    ADD CONSTRAINT collection_books_pkey PRIMARY KEY (collection_id, book_id);


--
-- Name: collections collections_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.collections
    ADD CONSTRAINT collections_pkey PRIMARY KEY (id);


--
-- Name: idx_books_author; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_books_author ON public.books USING btree (author);


--
-- Name: idx_books_genre; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_books_genre ON public.books USING btree (genre);


--
-- Name: idx_books_published_date; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_books_published_date ON public.books USING btree (published_date);


--
-- Name: collection_books collection_books_book_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.collection_books
    ADD CONSTRAINT collection_books_book_id_fkey FOREIGN KEY (book_id) REFERENCES public.books(id) ON DELETE CASCADE;


--
-- Name: collection_books collection_books_collection_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.collection_books
    ADD CONSTRAINT collection_books_collection_id_fkey FOREIGN KEY (collection_id) REFERENCES public.collections(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

