ALTER TABLE public."user" ADD verified bool DEFAULT false NOT NULL;

UPDATE public."user" set
    verified = true;
