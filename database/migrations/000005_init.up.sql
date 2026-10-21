CREATE TABLE IF NOT EXISTS public."transaction_split" (
  id BIGSERIAL not null,
  description character varying COLLATE pg_catalog."default" NOT NULL,
  detail character varying COLLATE pg_catalog."default",
  user_id bigint NOT NULL,
  value bigint NOT NULL,
  category_id bigint NOT NULL,
  transaction_id BIGINT not null,
  CONSTRAINT transaction_split_pk PRIMARY KEY (id, user_id),
  CONSTRAINT transaction_split_user_fk FOREIGN KEY (user_id)
        REFERENCES public."user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO action,
  CONSTRAINT transaction_split_category_id_fk FOREIGN KEY (category_id, user_id)
        REFERENCES public.category (id, user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO action,
  CONSTRAINT transaction_split_transaction_id_fk FOREIGN KEY (transaction_id, user_id)
        REFERENCES public.transaction (id, user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

insert into public."transaction_split" (description, user_id, value, category_id, transaction_id)
select description, user_id, value, category_id, id from public."transaction" t;

alter table "transaction" drop column description;
alter table "transaction" drop column detail;
alter table "transaction" drop column category_id;

alter table "transaction_tag" add column transaction_split_id bigint not null;

update "transaction_tag" set
  transaction_split_id = transaction_id;


alter table "transaction_tag" drop constraint transaction_tag_pk;
alter table "transaction_tag" drop constraint transaction_tag_transaction_fk;
alter table "transaction_tag" drop column transaction_id;
alter table "transaction_tag" add CONSTRAINT transaction_tag_pk PRIMARY KEY (transaction_split_id, tag, user_id);
alter table "transaction_tag" add CONSTRAINT transaction_tag_transaction_split_fk FOREIGN KEY (transaction_split_id, user_id) REFERENCES public."transaction_split"(id, user_id) ON DELETE CASCADE;

