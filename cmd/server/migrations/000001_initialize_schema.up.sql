CREATE SEQUENCE counters_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;
CREATE TABLE "public"."counters" (
    "id" integer DEFAULT nextval('counters_id_seq') NOT NULL,
    "name" character(40) NOT NULL UNIQUE,
    "value" integer NOT NULL,
    CONSTRAINT "counters_pkey" PRIMARY KEY ("id")
) WITH (oids = false);



CREATE SEQUENCE gauges_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;
CREATE TABLE "public"."gauges" (
    "id" integer DEFAULT nextval('gauges_id_seq') NOT NULL,
    "name" character(40) NOT NULL UNIQUE,
    "value" double precision NOT NULL,
    CONSTRAINT "gauges_pkey" PRIMARY KEY ("id")
) WITH (oids = false);
