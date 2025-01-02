-- Table: public.user

-- DROP TABLE IF EXISTS public."user";

CREATE TABLE IF NOT EXISTS public."user"
(
    id bigint NOT NULL DEFAULT nextval('user_id_seq'::regclass),
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

INSERT INTO `users` (`name`, `email`, `password`, `created_at`) VALUES ('Master', 'master@test.com', '"$2a$10$Cl84JkHPTsM0lWdgDOQxduCU6YRC9qR7glgDVIfJU0W2995sXLAGS"', curdate());

-- Table: public.category

-- DROP TABLE IF EXISTS public.category;

CREATE TABLE IF NOT EXISTS public.category
(
    id integer NOT NULL DEFAULT nextval('category_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    type character(1) COLLATE pg_catalog."default" NOT NULL,
    user_id bigint,
    CONSTRAINT category_pk PRIMARY KEY (id),
    CONSTRAINT category_user_fk FOREIGN KEY (user_id)
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
    id bigint NOT NULL DEFAULT nextval('account_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    type character(1) COLLATE pg_catalog."default" NOT NULL,
    due_day integer,
    user_id bigint NOT NULL,
    CONSTRAINT account_pk PRIMARY KEY (id),
    CONSTRAINT account_user_fk FOREIGN KEY (user_id)
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
    id bigint NOT NULL DEFAULT nextval('transaction_id_seq'::regclass),
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
    id bigint NOT NULL DEFAULT nextval('plan_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT plan_pk PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.plan
    OWNER to master;