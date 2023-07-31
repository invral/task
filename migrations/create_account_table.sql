BEGIN;

CREATE TABLE IF NOT EXISTS public.user (
                                           id SERIAL PRIMARY KEY NOT NULL,
                                           currency VARCHAR(3) NOT NULL,
                                           balance INT NOT NULL DEFAULT 0,
                                            password VARCHAR(255) NOT NULL,
                                            email VARCHAR(255) NOT NULL,
);

END;