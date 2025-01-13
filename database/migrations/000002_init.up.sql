-- Table: public.transaction_tag
ALTER TABLE public.transaction DROP CONSTRAINT transaction_pk;

ALTER TABLE public.transaction ADD CONSTRAINT transaction_pk PRIMARY KEY (id, user_id);

CREATE TABLE public.transaction_tag (
	transaction_id bigint NOT NULL,
	tag varchar NOT NULL,
	user_id bigint NOT NULL,
	CONSTRAINT transaction_tag_pk PRIMARY KEY (transaction_id, tag, user_id),
	CONSTRAINT transaction_tag_transaction_fk FOREIGN KEY (transaction_id,user_id) REFERENCES public."transaction"(id,id_user) ON DELETE CASCADE
);
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.transaction_tag
    OWNER to master;

