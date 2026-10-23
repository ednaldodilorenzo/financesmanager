--liquibase formatted sql

--changeset edsf:0001-baseline
--preconditions onFail:MARK_RAN
--precondition-sql-check expectedResult:0 SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'transaction'

CREATE TABLE public."user" (
                               id bigint NOT NULL,
                               name character varying,
                               email character varying NOT NULL,
                               password character varying,
                               created_at timestamp without time zone NOT NULL,
                               updated_at timestamp without time zone
);
CREATE SEQUENCE public.user_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER SEQUENCE public.user_id_seq OWNED BY public."user".id;
ALTER TABLE ONLY public."user" ALTER COLUMN id SET DEFAULT nextval('public.user_id_seq'::regclass);
ALTER TABLE ONLY public."user" ADD CONSTRAINT user_pk PRIMARY KEY (id);

CREATE TABLE public.account (
                                id bigint NOT NULL,
                                name character varying NOT NULL,
                                type character(1) NOT NULL,
                                due_day integer,
                                user_id bigint NOT NULL
);
CREATE SEQUENCE public.account_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER SEQUENCE public.account_id_seq OWNED BY public.account.id;
ALTER TABLE ONLY public.account ALTER COLUMN id SET DEFAULT nextval('public.account_id_seq'::regclass);
ALTER TABLE ONLY public.account ADD CONSTRAINT account_pk PRIMARY KEY (id, user_id);

CREATE TABLE public.category (
                                 id bigint NOT NULL,
                                 name character varying NOT NULL,
                                 type character(1) NOT NULL,
                                 user_id bigint NOT NULL,
                                 filter character varying NOT NULL
);
CREATE SEQUENCE public.category_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER SEQUENCE public.category_id_seq OWNED BY public.category.id;
ALTER TABLE ONLY public.category ALTER COLUMN id SET DEFAULT nextval('public.category_id_seq'::regclass);
ALTER TABLE ONLY public.category ADD CONSTRAINT category_pk PRIMARY KEY (id, user_id);

CREATE TABLE public.budget (
                               id bigint NOT NULL,
                               user_id bigint NOT NULL,
                               year integer NOT NULL,
                               category_id bigint NOT NULL,
                               value bigint NOT NULL
);
CREATE SEQUENCE public.budget_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER SEQUENCE public.budget_id_seq OWNED BY public.budget.id;
ALTER TABLE ONLY public.budget ALTER COLUMN id SET DEFAULT nextval('public.budget_id_seq'::regclass);
ALTER TABLE ONLY public.budget ADD CONSTRAINT budget_pk PRIMARY KEY (id, user_id);

CREATE TABLE public.transaction (
                                    id bigint NOT NULL,
                                    payment_date timestamp without time zone NOT NULL,
                                    transaction_date timestamp without time zone NOT NULL,
                                    value bigint NOT NULL,
                                    account_id bigint NOT NULL,
                                    user_id bigint NOT NULL
);
CREATE SEQUENCE public.transaction_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER SEQUENCE public.transaction_id_seq OWNED BY public.transaction.id;
ALTER TABLE ONLY public.transaction ALTER COLUMN id SET DEFAULT nextval('public.transaction_id_seq'::regclass);
ALTER TABLE ONLY public.transaction ADD CONSTRAINT transaction_pk PRIMARY KEY (id, user_id);

CREATE TABLE public.transaction_split (
                                          id bigint NOT NULL,
                                          description character varying NOT NULL,
                                          detail character varying,
                                          user_id bigint NOT NULL,
                                          value bigint NOT NULL,
                                          category_id bigint NOT NULL,
                                          transaction_id bigint NOT NULL
);
CREATE SEQUENCE public.transaction_split_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER SEQUENCE public.transaction_split_id_seq OWNED BY public.transaction_split.id;
ALTER TABLE ONLY public.transaction_split ALTER COLUMN id SET DEFAULT nextval('public.transaction_split_id_seq'::regclass);
ALTER TABLE ONLY public.transaction_split ADD CONSTRAINT transaction_split_pk PRIMARY KEY (id, user_id);

CREATE TABLE public.transaction_tag (
                                        tag character varying NOT NULL,
                                        user_id bigint NOT NULL,
                                        transaction_split_id bigint DEFAULT 1 NOT NULL
);
ALTER TABLE ONLY public.transaction_tag ADD CONSTRAINT transaction_tag_pk PRIMARY KEY (transaction_split_id, tag, user_id);

ALTER TABLE ONLY public.account ADD CONSTRAINT account_user_fk FOREIGN KEY (user_id) REFERENCES public."user"(id);
ALTER TABLE ONLY public.category ADD CONSTRAINT category_user_fk FOREIGN KEY (user_id) REFERENCES public."user"(id);
ALTER TABLE ONLY public.budget ADD CONSTRAINT budget_category_fk FOREIGN KEY (category_id, user_id) REFERENCES public.category(id, user_id);
ALTER TABLE ONLY public.transaction ADD CONSTRAINT transaction_user_fk FOREIGN KEY (user_id) REFERENCES public."user"(id);
ALTER TABLE ONLY public.transaction ADD CONSTRAINT transaction_account_id_fk FOREIGN KEY (account_id, user_id) REFERENCES public.account(id, user_id);
ALTER TABLE ONLY public.transaction_split ADD CONSTRAINT transaction_split_user_fk FOREIGN KEY (user_id) REFERENCES public."user"(id);
ALTER TABLE ONLY public.transaction_split ADD CONSTRAINT transaction_split_category_id_fk FOREIGN KEY (category_id, user_id) REFERENCES public.category(id, user_id);
ALTER TABLE ONLY public.transaction_split ADD CONSTRAINT transaction_split_transaction_id_fk FOREIGN KEY (transaction_id, user_id) REFERENCES public.transaction(id, user_id);
ALTER TABLE ONLY public.transaction_tag ADD CONSTRAINT transaction_tag_transaction_split_fk FOREIGN KEY (transaction_split_id, user_id) REFERENCES public.transaction_split(id, user_id) ON DELETE CASCADE;
