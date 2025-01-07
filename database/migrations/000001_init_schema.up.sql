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
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public."user"
    OWNER to master;

insert into public.user (name, email, password, created_at) values ('Ednaldo', 'ednaldo.dilorenzo@gmail.com', '$2a$10$I2bEvoneiS72A0wGgkDfNei0M1QasNuJFWp.NvlfNbBsnSVYSBpJG', '2024-11-02');

-- Table: public.category

-- DROP TABLE IF EXISTS public.category;

CREATE TABLE IF NOT EXISTS public.category
(
    id bigserial NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    type character(1) COLLATE pg_catalog."default" NOT NULL,
    id_user bigint,
    CONSTRAINT category_pk PRIMARY KEY (id),
    CONSTRAINT category_user_fk FOREIGN KEY (id_user)
        REFERENCES public."user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.category
    OWNER to master;

  -- Table: public.account

-- DROP TABLE IF EXISTS public.account;

CREATE TABLE IF NOT EXISTS public.account
(
    id bigserial NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    type character(1) COLLATE pg_catalog."default" NOT NULL,
    due_day integer,
    id_user bigint NOT NULL,
    CONSTRAINT account_pk PRIMARY KEY (id),
    CONSTRAINT account_user_fk FOREIGN KEY (id_user)
        REFERENCES public."user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.account
    OWNER to master;

-- Table: public.transaction

-- DROP TABLE IF EXISTS public.transaction;

CREATE TABLE IF NOT EXISTS public.transaction
(
    id bigserial NOT NULL,
    description character varying COLLATE pg_catalog."default" NOT NULL,
    detail character varying COLLATE pg_catalog."default",
    payment_date timestamp without time zone NOT NULL,
    transaction_date timestamp without time zone NOT NULL,
    value bigint NOT NULL,
    id_category bigint NOT NULL,
    id_account bigint NOT NULL,
    id_user bigint NOT NULL,
    CONSTRAINT transaction_pk PRIMARY KEY (id),
    CONSTRAINT transaction_account_id_fk FOREIGN KEY (id_account)
        REFERENCES public.account (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT transaction_category_id_fk FOREIGN KEY (id_category)
        REFERENCES public.category (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT transaction_user_fk FOREIGN KEY (id_user)
        REFERENCES public."user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.transaction
    OWNER to master;

-- Table: public.plan

-- DROP TABLE IF EXISTS public.plan;

CREATE TABLE IF NOT EXISTS public.plan
(
    id bigserial NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT plan_pk PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.plan
    OWNER to master;
