CREATE TABLE public.budget (
        id bigserial NOT NULL,
	user_id bigint NOT NULL,
	"year" int NOT NULL,
	category_id bigint NOT NULL,
	value int8 NOT NULL,
	CONSTRAINT budget_pk PRIMARY KEY (id, user_id),
	CONSTRAINT budget_category_fk FOREIGN KEY (category_id,user_id) REFERENCES public.category(id,user_id)
);

