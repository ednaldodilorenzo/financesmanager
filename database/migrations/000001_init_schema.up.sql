-- Table: public.user

-- DROP TABLE IF EXISTS public."user";

CREATE TABLE IF NOT EXISTS public."user"
(
    id bigserial NOT NULL,
    name character varying COLLATE pg_catalog."default",
    email character varying COLLATE pg_catalog."default" NOT NULL,
    password character varying COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone,
    CONSTRAINT user_pk PRIMARY KEY (id)
);

insert into public.user (name, email, password, created_at) values ('Ednaldo', 'ednaldo.dilorenzo@gmail.com', '$2a$10$I2bEvoneiS72A0wGgkDfNei0M1QasNuJFWp.NvlfNbBsnSVYSBpJG', '2024-11-02');

-- Table: public.category

-- DROP TABLE IF EXISTS public.category;

CREATE TABLE IF NOT EXISTS public.category
(
    id bigserial NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    type character(1) COLLATE pg_catalog."default" NOT NULL,
    user_id bigint,
    filter varchar NOT NULL,
    CONSTRAINT category_pk PRIMARY KEY (id, user_id),
    CONSTRAINT category_user_fk FOREIGN KEY (user_id)
        REFERENCES public."user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

  -- Table: public.account

-- DROP TABLE IF EXISTS public.account;

CREATE TABLE IF NOT EXISTS public.account
(
    id bigserial NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    type character(1) COLLATE pg_catalog."default" NOT NULL,
    due_day integer,
    user_id bigint NOT NULL,
    CONSTRAINT account_pk PRIMARY KEY (id, user_id),
    CONSTRAINT account_user_fk FOREIGN KEY (user_id)
        REFERENCES public."user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);


-- Table: public.transaction

-- DROP TABLE IF EXISTS public.transaction;

CREATE TABLE IF NOT EXISTS public."transaction"
(
    id bigserial NOT NULL,
    description character varying COLLATE pg_catalog."default" NOT NULL,
    detail character varying COLLATE pg_catalog."default",
    payment_date timestamp without time zone NOT NULL,
    transaction_date timestamp without time zone NOT NULL,
    value bigint NOT NULL,
    category_id bigint NOT NULL,
    account_id bigint NOT NULL,
    user_id bigint NOT NULL,
    CONSTRAINT transaction_pk PRIMARY KEY (id, user_id),
    CONSTRAINT transaction_account_id_fk FOREIGN KEY (account_id, user_id)
        REFERENCES public.account (id, user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT transaction_category_id_fk FOREIGN KEY (category_id, user_id)
        REFERENCES public.category (id, user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT transaction_user_fk FOREIGN KEY (user_id)
        REFERENCES public."user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

-- Table:

CREATE TABLE public.transaction_tag (
	transaction_id bigint NOT NULL,
	tag varchar NOT NULL,
	user_id bigint NOT NULL,
	CONSTRAINT transaction_tag_pk PRIMARY KEY (transaction_id, tag, user_id),
	CONSTRAINT transaction_tag_transaction_fk FOREIGN KEY (transaction_id, user_id) REFERENCES public."transaction"(id, user_id) ON DELETE CASCADE
);
